package jirachat

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	issueLinkBase = "https://%s.atlassian.net/browse/%s"
	userLinkBase  = "https://%s.atlassian.net/secure/ViewProfile.jspa?name=%s"
)

var ErrSlackParse = errors.New("unknown Event Failed Slack Parsing")

// SlackMessage represents a payload sent to Slack.
// The values are sent to Slack via incoming-webhook.
// See - https://my.slack.com/services/new/incoming-webhook
type SlackMessage struct {
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

// Default constructSlackMessage for issue_updated type. Unfortunately this includes
// everything that isn't worklog or ticket create/delete
func (s *SlackService) IssueUpdated(event *JIRAWebevent) error {

	payload := SlackMessage{}
	var fields []Field
	title := ""
	user := event.GetUserLink(s.Config)
	// Try to determine what kind of event this was
	switch {
	case len(event.Comment.Id) > 0:
		title = fmt.Sprintf("%s commented on %s", user,
			event.GetIssueLink(s.Config))
		fields = []Field{
			{
				Title: "Issue",
				Value: event.Issue.Fields.Summary,
				Short: false,
			},
			{
				Title: "Comment",
				Value: event.Comment.Body,
				Short: false,
			},
		}
	case len(event.Changelog.Items) > 0:
		switch {
		case event.Changelog.Items[0].Field == "status":
			title = fmt.Sprintf("%s changed status of %s", user,
				event.GetIssueLink(s.Config))
			fields = []Field{
				{
					Title: "From",
					Value: event.Changelog.Items[0].FromString,
					Short: false,
				},
				{
					Title: "To",
					Value: event.Changelog.Items[0].ToString,
					Short: false,
				},
			}
		case event.Changelog.Items[0].Field == "assignee":
			title = fmt.Sprintf("%s changed assigne of %s", user,
				event.GetIssueLink(s.Config))

			from := "unassigned"
			if len(event.Changelog.Items[0].FromString) > 0 {
				from = event.Changelog.Items[0].FromString
			}
			to := "unassigned"
			if len(event.Changelog.Items[0].ToString) > 0 {
				to = event.Changelog.Items[0].ToString
			}
			fields = []Field{
				{
					Title: "From",
					Value: from,
					Short: false,
				},
				{
					Title: "To",
					Value: to,
					Short: false,
				},
			}
		default:
			// Post a generic event and post the details to the error channel
			title = fmt.Sprintf("%s modified %s", event.User.DisplayName,
				event.GetIssueLink(s.Config))
			resp := &Response{"Erroring Event": event}
			SendErrorNotice(resp.String(), s.Config)
			return ErrSlackParse

		}
	default:
		// Post a generic event and post the details to the error channel
		title = fmt.Sprintf("%s modified %s", event.User.DisplayName,
			event.GetIssueLink(s.Config))
		resp := &Response{"Erroring Event": event}
		SendErrorNotice(resp.String(), s.Config)
		return ErrSlackParse
	}

	attachment := Attachment{
		Fallback: title,
		Pretext:  title,
		Color:    event.GetPriorityColor(),
		Fields:   fields,
	}

	payload.Channel = s.Config.Channel
	payload.Username = s.Config.BotName
	payload.Icon_url = event.User.LargeAvatar()
	payload.Unfurl_links = true
	payload.Text = ""
	payload.Attachments = []Attachment{attachment}
	return payload.SendEvent(s.Config)
}

// Default construct SlackMessage for issue_created type
func (s *SlackService) IssueCreated(event *JIRAWebevent) error {
	payload := SlackMessage{}
	fields := []Field{
		{
			Title: "Summary",
			Value: event.Issue.Fields.Summary,
			Short: false,
		},
		{
			Title: "Assignee",
			Value: event.Issue.Fields.Assignee.DisplayName,
			Short: true,
		},
		{
			Title: "Priority",
			Value: event.Issue.Fields.Priority.Name,
			Short: true,
		},
	}
	title := fmt.Sprintf("%s created %s", event.GetUserLink(s.Config),
		event.GetIssueLink(s.Config))
	attachment := Attachment{
		Fallback: title,
		Pretext:  title,
		Color:    event.GetPriorityColor(),
		Fields:   fields,
	}

	payload.Channel = s.Config.Channel
	payload.Username = s.Config.BotName
	payload.Icon_url = event.User.LargeAvatar()
	payload.Unfurl_links = true
	payload.Text = ""
	payload.Attachments = []Attachment{attachment}
	return payload.SendEvent(s.Config)
}

// Default construct SlackMessage for issue_deleted type
func (s *SlackService) IssueDeleted(event *JIRAWebevent) error {
	payload := SlackMessage{}
	body := "None"
	last := event.Issue.Fields.Comment.Total
	if last > 0 {
		body = event.Issue.Fields.Comment.Comments[last-1].Body
	}

	fields := []Field{
		{
			Title: "Issue",
			Value: event.Issue.Fields.Summary,
			Short: false,
		},
		{
			Title: "Last Comment",
			Value: body,
			Short: false,
		},
	}

	// Don't bother linking to the issue!
	title := fmt.Sprintf("%s deleted %s", event.GetUserLink(s.Config),
		event.Issue.Key)
	attachment := Attachment{
		Fallback: title,
		Pretext:  title,
		Fields:   fields,
	}

	payload.Channel = s.Config.Channel
	payload.Username = s.Config.BotName
	payload.Icon_url = event.User.LargeAvatar()
	payload.Unfurl_links = true
	payload.Text = ""
	payload.Attachments = []Attachment{attachment}
	return payload.SendEvent(s.Config)
}

// Default construct SlackMessage for issue_deleted type
func (s *SlackService) WorklogUpdated(event *JIRAWebevent) error {
	payload := SlackMessage{}

	timestr := ""
	for i := range event.Changelog.Items {
		if event.Changelog.Items[i].Field == "timespent" {
			timestr = event.Changelog.Items[i].ToString
		}
	}
	if len(timestr) == 0 {
		return errors.New("Unable to read timespent field")
	}

	time, err := strconv.Atoi(timestr)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid timespent field %s", timestr))
	}
	time /= 60

	if time == 1 {
		timestr = strconv.Itoa(time) + " minute"
	} else {
		timestr = strconv.Itoa(time) + " minutes"
	}

	fields := []Field{
		{
			Title: "Total Work",
			Value: timestr,
			Short: false,
		},
	}

	title := fmt.Sprintf("%s updated work log %s", event.GetUserLink(s.Config),
		event.GetIssueLink(s.Config))
	attachment := Attachment{
		Fallback: title,
		Pretext:  title,
		Color:    event.GetPriorityColor(),
		Fields:   fields,
	}

	payload.Channel = s.Config.Channel
	payload.Username = s.Config.BotName
	payload.Icon_url = event.User.LargeAvatar()
	payload.Unfurl_links = true
	payload.Text = ""
	payload.Attachments = []Attachment{attachment}
	return payload.SendEvent(s.Config)
}

func(s *SlackService) CommentCreated(event *JIRAWebevent) error{
	payload := SlackMessage{}
	var fields []Field
	title := ""
	user := event.GetUserLink(s.Config)
	title = fmt.Sprintf("%s commented on %s", user,
		event.GetIssueLink(s.Config))
	fields = []Field{
		{
			Title: "Issue",
			Value: event.Issue.Fields.Summary,
			Short: false,
		},
		{
			Title: "Comment",
			Value: event.Comment.Body,
			Short: false,
		}}

	attachment := Attachment{
		Fallback: title,
		Pretext:  title,
		Color:    event.GetPriorityColor(),
		Fields:   fields,
	}

	payload.Channel = s.Config.Channel
	payload.Username = s.Config.BotName
	payload.Icon_url = event.User.LargeAvatar()
	payload.Unfurl_links = true
	payload.Text = ""
	payload.Attachments = []Attachment{attachment}
	return payload.SendEvent(s.Config)
}

// Returns a markdown formatted issue link with the issue key
// as the link text
func (e *JIRAWebevent) GetIssueLink(s *SlackConfig) string {
	link := fmt.Sprintf(issueLinkBase, s.Domain, e.Issue.Key)
	return fmt.Sprintf("<%s|%s>", link, e.Issue.Key)
}

// Returns a markdown formatted user link with the user name
// as the link text
func (e *JIRAWebevent) GetUserLink(s *SlackConfig) string {
	link := fmt.Sprintf(userLinkBase, s.Domain, e.User.Name)
	return fmt.Sprintf("<%s|%s>", link, e.User.DisplayName)
}

// Convert priority id to hex color string
func (e *JIRAWebevent) GetPriorityColor() string {

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
