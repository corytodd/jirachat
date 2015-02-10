package jirachat

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const jira_img = "https://dujrsrsgsd3nh.cloudfront.net/img/emoticons/jira-1350074257.png"

// This is a json response for a JIRA webhook (more or less) according to
//https://developer.atlassian.com/display/JIRADEV/JIRA+Webhooks+Overview
type JIRAWebevent struct {
	// Internal ID of the event
	Id int `json:"id,omitempty"`

	// Unix timestamp
	Timestamp int `json:"timestamp,omitempty"`

	// Object describing issue event relates to
	Issue JIRAIssue `json:"issue"`

	// User who triggered event
	User JIRAUser `json:"user"`

	Changelog JIRAChangelog `json:"changelog"`

	// Set if this event is a jira_updated event and a comment was made
	Comment JIRAComment `json:"comment"`

	// The type of event
	WebhookEvent string `json:"webhookEvent"`
}

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
	FieldType  string `json:fieldtype"`
	Field      string `json:"field"`
}

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
}

type JIRAIssueAssignee struct {
	Self        string            `json:"self"`
	Name        string            `json:"name"`
	Key         string            `json:"key"`
	Email       string            `json:"emailAddress"`
	AvatarUrls  map[string]string `json:"avatarUrls"`
	DisplayName string            `json:"displayName"`
	Active      bool              `json:"active"`
	timeZone    string            `json:"timeZone"`
}

type JIRAIssuePriority struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type JIRAIssueStatus struct {
	Name string `json:"name"`
}

type InnerComment struct {
	StartAt    int           `json"startAt"`
	MaxResults int           `json:"maxResults"`
	Total      int           `json:"total"`
	Comments   []JIRAComment `json:"comments"`
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
