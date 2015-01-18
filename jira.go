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
	Id           int           `json:"id,omitempty"`
	Timestamp    int           `json:"timestamp,omitempty"`
	Issue        JIRAIssue     `json:"issue"`
	User         JIRAUser      `json:"user"`
	Changelog    JIRAChangelog `json:"changelog"`
	Comment      JIRAComment   `json:"comment"`
	WebhookEvent string        `json:"webhookEvent"`
}

type JIRAIssue struct {
	Expand string         `json:"expand"`
	Id     string         `json:"id"`
	Self   string         `json:"self"`
	Key    string         `json:"key"`
	Fields IssueFieldData `json:"fields"`
}

type JIRAUser struct {
	Self         string            `json:"self"`
	Name         string            `json:"name"`
	EmailAddress string            `json:"emailAddress"`
	AvatarUrls   map[string]string `json:"avatarUrls"`
	DisplayName  string            `json:"displayName"`
	Active       bool              `json:"active"`
}

// Some of the ChangeLogItems may through unmarshal errors but they don't seem
// to cause any major issues
type JIRAChangelog struct {
	Items []ChangleLogItems `json:"items,omitempty"`
	Id    int               `json:"id,omitempty"`
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
	Summary     string `json:"summary"`
	Created     string `json:"created"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Assignee    string `json:"assignee"`
}

// For the avatars, use getter methods because the names start with numbers
func (j *JIRAUser) SmallAvatar() string {
	return j.AvatarUrls["16x16"]
}
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
	// is is safe to ignore.
	err = json.Unmarshal(body, &event)
	return event, err
}
