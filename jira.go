package jirachat

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const jira_img = "https://dujrsrsgsd3nh.cloudfront.net/img/emoticons/jira-1350074257.png"

// This is a json response for a JIRA webhook (more or less) according to
//https://developer.atlassian.com/display/JIRADEV/JIRA+Webhooks+Overview
type JiraWebevent struct {
	Id           int           `json:"id,omitempty"`
	Timestamp    string        `json:"timestamp,omitempty"`
	Issue        JiraIssue     `json:"issue"`
	User         JiraUser      `json:"user"`
	Changelog    JiraChangelog `json:"changelog"`
	Comment      JiraComment   `json:"comment"`
	WebhookEvent string        `json:"webhookEvent"`
}

type JiraIssue struct {
	Expand string    `json:"expand"`
	Id     string    `json:"id"`
	Self   string    `json:"self"`
	Key    string    `json:"key"`
	Fields FieldData `json:"fields"`
}

type JiraUser struct {
	Self         string            `json:"self"`
	Name         string            `json:"name"`
	EmailAddress string            `json:"emailAddress"`
	AvatarUrls   map[string]string `json:"avatarUrls"`
	DisplayName  string            `json:"displayName"`
	Active       string            `json:"active"`
}

type JiraChangelog struct {
	Items map[string]interface{} `json:"items"`
	Id    int                    `json:"id,omitempty"`
}

type JiraComment struct {
	Self         string   `json:"self"`
	Id           string   `json:"id"`
	Author       JiraUser `json:"author"`
	Body         string   `json:"body"`
	UpdateAuthor JiraUser `json:"updateAuthor"`
	Created      string   `json:"created"`
	Updated      string   `json:"updated"`
}

type FieldData struct {
	Summary     string `json:"summary"`
	Created     string `json:"created"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Assignee    string `json:"assignee"`
}

func (j *JiraUser) SmallAvatar() string {
	return j.AvatarUrls["16x16"]
}
func (j *JiraUser) LargeAvatar() string {
	return j.AvatarUrls["48x48"]
}

//TODO add JiraIssue.Fields.labels function(s)

func (f *JiraChangelog) GetToString() (string, error) {
	result, ok := f.Items["toString"].(string)
	if !ok {
		return "", errors.New("jirachangelog: No toString")
	}
	return result, nil
}
func (f *JiraChangelog) GetTo() (string, error) {
	result, ok := f.Items["to"].(string)
	if !ok {
		return "", errors.New("jirachangelog: No to")
	}
	return result, nil
}
func (f *JiraChangelog) GetFromString() (string, error) {
	result, ok := f.Items["fromString"].(string)
	if !ok {
		return "", errors.New("jirachangelog: No fromString")
	}
	return result, nil
}
func (f *JiraChangelog) GetFrom() (string, error) {
	result, ok := f.Items["from"].(string)
	if !ok {
		return "", errors.New("jirachangelog: No from")
	}
	return result, nil
}
func (f *JiraChangelog) GetFieldType() (string, error) {
	result, ok := f.Items["fieldtype"].(string)
	if !ok {
		return "", errors.New("jirachangelog: No fieldtype")
	}
	return result, nil
}
func (f *JiraChangelog) GetField() (string, error) {
	result, ok := f.Items["field"].(string)
	if !ok {
		return "", errors.New("jirachangelog: No field")
	}
	return result, nil
}

// Parse the request body as a JIRA webhook event
// Returns a new JiraWebEvent object or error
func Parse(r *http.Request) (JiraWebevent, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var event JiraWebevent

	// This will generate a error unmarshaling some of the data but
	// is is safe to ignore.
	err = json.Unmarshal(body, &event)
	return event, err
}
