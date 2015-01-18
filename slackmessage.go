package jirachat

import (
	"fmt"
)

// Payload represents a payload sent to Slack.
// The values are sent to Slack via incoming-webhook.
// See - https://my.slack.com/services/new/incoming-webhook
type Payload struct {
	Channel      string       `json:"channel"`
	Username     string       `json:"username"`
	Text         string       `json:"text"`
	Icon_emoji   string       `json:"icon_emoji"`
	Icon_url     string       `json:"icon_url"`
	Unfurl_links bool         `json:"unfurl_links"`
	Attachments  []Attachment `json:"attachments"`
}

// Attachment is an attachment to Payload.
// The format is defined in Slack Api document.
// See - https://api.slack.com/docs/attachments
type Attachment struct {

	// text summary of the attachment that is shown by clients that understand
	// attachments but choose not to show them.
	// Required
	Fallback string `json:"fallback"`

	// text that should appear within the attachment
	// Optional
	Text string `json:"text"`

	// text that should appear above the formatted data",
	// Optional
	Pretext string `json:"pretext"`

	// Can either be one of 'good', 'warning', 'danger', or any hex color code
	// Optional
	Color string `json:"color"`

	// Fields are displayed in a table on the message
	Fields []Field `json:"fields"`
}

// Field is a field to Attachment.
// Like Attachment, the format is defined in Slack Api document.
// see - https://api.slack.com/docs/attachments
type Field struct {
	// The title may not contain markup and will be escaped for you
	// Required
	Title string `json:"title"`

	// Text value of the field. May contain standard message markup and
	// must be escaped as normal. May be multi-line.",
	Value string `json:"value"`

	// flag indicating whether the `value` is short enough to be
	// displayed side-by-side with other values
	// Optional
	Short bool `json:"short"`
}

// ConstructSlackMessage for issue_updated type
func (s *slackService) IssueUpdated(event JIRAWebevent) error {
	payload := Payload{}
	fields := []Field{
		Field{
			Title: "Issue",
			Value: event.Issue.Fields.Summary,
			Short: false,
		},
		Field{
			Title: "Comment",
			Value: event.Comment.Body,
			Short: false,
		},
	}

	attachment := Attachment{
		Fallback: fmt.Sprintf("%s Commented on <%s|%s>", event.User.DisplayName,
			event.Issue.Self, event.Issue.Key),
		Pretext: fmt.Sprintf("%s Commented on <%s|%s>", event.User.DisplayName,
			event.Issue.Self, event.Issue.Key),
		Color:  event.getPriorityColor(),
		Fields: fields,
	}

	payload.Channel = s.config_.Channel
	payload.Username = s.config_.BotName
	payload.Icon_url = event.User.LargeAvatar()
	payload.Unfurl_links = true
	payload.Text = ""
	payload.Attachments = []Attachment{attachment}
	return payload.sendEvent(s.config_)
}

// ConstructSlackMessage for issue_created type
func (s *slackService) IssueCreated(event JIRAWebevent) error {
	payload := Payload{}
	fields := []Field{
		Field{
			Title: "Summary",
			Value: event.Issue.Fields.Summary,
			Short: false,
		},
		Field{
			Title: "Assignee",
			Value: event.Issue.Fields.Assignee.DisplayName,
			Short: true,
		},
		Field{
			Title: "Priority",
			Value: event.Issue.Fields.Priority.Name,
			Short: true,
		},
	}

	attachment := Attachment{
		Fallback: fmt.Sprintf("%s Created on <%s|%s>", event.User.DisplayName,
			event.Issue.Self, event.Issue.Key),
		Pretext: fmt.Sprintf("%s Created on <%s|%s>", event.User.DisplayName,
			event.Issue.Self, event.Issue.Key),
		Color:  event.getPriorityColor(),
		Fields: fields,
	}

	payload.Channel = s.config_.Channel
	payload.Username = s.config_.BotName
	payload.Icon_url = event.User.LargeAvatar()
	payload.Unfurl_links = true
	payload.Text = ""
	payload.Attachments = []Attachment{attachment}
	return payload.sendEvent(s.config_)
}

// ConstructSlackMessage for issue_deleted type
func (s *slackService) IssueDeleted(event JIRAWebevent) error {
	payload := Payload{}
	fields := []Field{
		Field{
			Title: "Issue",
			Value: event.Issue.Fields.Summary,
			Short: false,
		},
		Field{
			Title: "Comment",
			Value: event.Comment.Body,
			Short: false,
		},
	}

	attachment := Attachment{
		Fallback: fmt.Sprintf("%s Commented on <%s|%s>", event.User.DisplayName,
			event.Issue.Self, event.Issue.Key),
		Pretext: fmt.Sprintf("%s Commented on <%s|%s>", event.User.DisplayName,
			event.Issue.Self, event.Issue.Key),
		Color:  event.getPriorityColor(),
		Fields: fields,
	}

	payload.Channel = s.config_.Channel
	payload.Username = s.config_.BotName
	payload.Icon_url = event.User.LargeAvatar()
	payload.Unfurl_links = true
	payload.Text = ""
	payload.Attachments = []Attachment{attachment}
	return payload.sendEvent(s.config_)
}

// ConstructSlackMessage for worklog updates
func (s *slackService) WorklogUpdated(event JIRAWebevent) error {
	payload := Payload{}
	fields := []Field{
		Field{
			Title: "Issue",
			Value: event.Issue.Fields.Summary,
			Short: false,
		},
		Field{
			Title: "Comment",
			Value: event.Comment.Body,
			Short: false,
		},
	}

	attachment := Attachment{
		Fallback: fmt.Sprintf("%s Commented on <%s|%s>", event.User.DisplayName,
			event.Issue.Self, event.Issue.Key),
		Pretext: fmt.Sprintf("%s Commented on <%s|%s>", event.User.DisplayName,
			event.Issue.Self, event.Issue.Key),
		Color:  event.getPriorityColor(),
		Fields: fields,
	}

	payload.Channel = s.config_.Channel
	payload.Username = s.config_.BotName
	payload.Icon_url = event.User.LargeAvatar()
	payload.Unfurl_links = true
	payload.Text = ""
	payload.Attachments = []Attachment{attachment}
	return payload.sendEvent(s.config_)
}

func (e *JIRAWebevent) getPriorityColor() string {

	id := e.Issue.Fields.Priority.Id
	switch {
	case id == "1": // Blocker
		return "#990000"
	case id == "2":
		return "#cc0000" // Critical
	case id == "3":
		return "#ff0000"
	case id == "6": // Normal
		return "#339933"
	case id == "4": // Minor
		return "#006600"
	case id == "5": // Trivial
		return "#003300"
	case id == "10000": // Holding
		return "#000000"
	default:
		return "good"
	}
}
