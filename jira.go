package jirachat

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
	"strings"
)

// This is a json response for a JIRA webhook (more or less) according to
// https://developer.atlassian.com/jiradev/jira-architecture/webhooks
type JIRAWebevent struct {
	// Internal ID of the event
	Id int `json:"id,omitempty"`

	// Unix timestamp
	Timestamp int `json:"timestamp,omitempty"`

	// Object describing issue event relates to
	Issue JIRAIssue `json:"issue"`

	// User who triggered event
	User JIRAUser `json:"user"`

	// An array of change items, with one entry for each field that has
	// been changed. The changelog is only provided for the issue_updated event.
	Changelog JIRAChangelog `json:"changelog"`

	// Set if this event is a jira_updated event and a comment was made
	Comment JIRAComment `json:"comment"`

	// The type of event
	WebhookEvent string `json:"webhookEvent"`
}

// Decribes the JIRAIssue object defined in the JIRA 5.1 REST docs
//
// https://docs.atlassian.com/jira/REST/5.1/#id204637
type JIRAIssue struct {
	Expand string `json:"expand"`

	// Internal ID of the issue
	Id string `json:"id"`

	Self string `json:"self"`

	// This is the common issue key, e.g. JIRA-100
	Key string `json:"key"`

	// There are quite a few fields and that ones provided in this library
	// are nowhere near exhaustive.
	Fields IssueFieldData `json:"fields"`
}

// Describes the JIRAUser object defined in the JIRA 5.1 REST docs
//
// https://docs.atlassian.com/jira/REST/5.1/#id202197
type JIRAUser struct {
	Self string `json:"self"`

	// The user's system name, e.g. mmcfly
	Name         string            `json:"name"`
	EmailAddress string            `json:"emailAddress"`
	AvatarUrls   map[string]string `json:"avatarUrls"`

	// Pretty name E.g. Mart McFly
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
}

// Some of the ChangeLogItems may through unmarshal errors but they don't seem
// to cause any major issues
type JIRAChangelog struct {
	Items []ChangleLogItems `json:"items,omitempty"`
	Id    string            `json:"id,omitempty"`
}

type ChangleLogItems struct {
	ToString   string `json:"toString"`
	To         string `json:"to"`
	FromString string `json:"fromString"`
	From       string `json:"from"`
	FieldType  string `json:"fieldtype"`
	Field      string `json:"field"`
}

// Describes the JIRAComment object defined in the JIRA 5.1 REST docs
//
// https://docs.atlassian.com/jira/REST/5.1/#id204337
type JIRAComment struct {
	Self         string   `json:"self"`
	Id           string   `json:"id"`
	Author       JIRAUser `json:"author"`
	Body         string   `json:"body"`
	UpdateAuthor JIRAUser `json:"updateAuthor"`
	Created      string   `json:"created"`
	Updated      string   `json:"updated"`
}

type IssueFieldData struct {
	Summary     string            `json:"summary"`
	Created     string            `json:"created"`
	Description string            `json:"description"`
	Priority    JIRAIssuePriority `json:"priority"`
	Assignee    JIRAIssueAssignee `json:"assignee"`
	Labels      []string          `json:"labels"`
	Status      JIRAIssueStatus   `json:"status"`
	Comment     InnerComment      `json:"comment"`
	IssueType   JIRAIssueType     `json:"issuetype"`
	Project     JIRAProject       `json:"project"`
	// CustomFields is a map of customfield_xxx from your JIRA instance. The key will match whichever
	// custom fields you have created. The contents obviously depend on what you have created. The value
	// the raw string value of whatever your field contains.
	CustomFields map[string]string
}

type JIRAIssueAssignee struct {
	Self        string            `json:"self"`
	Name        string            `json:"name"`
	Key         string            `json:"key"`
	Email       string            `json:"emailAddress"`
	AvatarUrls  map[string]string `json:"avatarUrls"`
	DisplayName string            `json:"displayName"`
	Active      bool              `json:"active"`
	Timezone    string            `json:"timeZone"`
}

type JIRAIssuePriority struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type JIRAIssueStatus struct {
	Name string `json:"name"`
}

type JIRAIssueType struct {
	Self        string `json:"self"`
	Id          string `json:"id"`
	Description string `json:"description"`
	IconURL     string `json:"iconUrl"`
	Name        string `json:"name"`
	subtask     bool   `json:"subtask"`
}

type InnerComment struct {
	StartAt    int           `json:"startAt"`
	MaxResults int           `json:"maxResults"`
	Total      int           `json:"total"`
	Comments   []JIRAComment `json:"comments"`
}

type JIRAProject struct {
	Self       string            `json:"self"`
	Id         string            `json:"id"`
	Key        string            `json:"key"`
	Name       string            `json:"name"`
	IconUrl    string            `json:"iconUrl"`
	Subtask    bool              `json:"subtask"`
	AvatarUrls map[string]string `json:"avatarUrls"`
}

// Returns the 16x16 user Avatar
// Note: If this is not a Gravatar, it will not render
// in Slack message because uploaded images are private to your JIRA instance.
func (j *JIRAUser) SmallAvatar() string {
	return j.AvatarUrls["16x16"]
}

// Returns the 48x48 user Avatar
// Note: If this is not a Gravatar, it will not render
// in Slack message because uploaded images are private to your JIRA instance.
func (j *JIRAUser) LargeAvatar() string {
	return j.AvatarUrls["48x48"]
}

//TODO add JIRAIssue.Fields.labels function(s)

// Parse the request body as a JIRA webhook event
// Returns a new JiraWebEvent object or error
func Parse(r *http.Request) (JIRAWebevent, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var event JIRAWebevent

	// This will generate a error unmarshaling some of the data but
	// is is safe to ignore. We return the error so you at least know
	// that there is some oddly formed data.
	err = json.Unmarshal(body, &event)

	event.Issue.Fields.CustomFields = make(map[string]string, 0)

	if len(event.Issue.Id) != 0 {
		jsonparser.ObjectEach(body, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			k := string(key)
			if strings.HasPrefix(k, "customfield_") {
				event.Issue.Fields.CustomFields[string(key)] = string(value)
			}
			return nil
		}, "issue", "fields")
	}

	return event, err
}

// Convenience interface for printing anonymous JSON objects
type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}
