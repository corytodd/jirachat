package jirachat

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type ChannelService struct {
	client *SlackClient
}

type SlackMessage struct {
	// Authentication token (Requires scope: post)
	// Format as xxxx-xxxxxxxxx-xxxx
	// Required
	Token string `json:"token"`

	// Channel to send message to. Can be a public channel, private group or IM
	// channel. Can be an encoded ID, or a name.
	// Required
	Channel string `json"channel"`

	// Text of the message to send
	// Required
	Text string `json:"text"`

	// From field in Slack chat
	// Optional
	BotName string `json:"username,omitempty"`

	// Change how messages are treated. See below.
	// Optional
	Parse string `json:"parse,omitempty"`

	// Find and link channel names and usernames.
	// Optional
	Link_names int `json:"link_names,omitempty"`

	// Structured message attachments
	// Optional
	Attachment Attachment `json:"fields,omitempty"`

	// Optional
	Unfurl_links bool `json:"unfurl_links,omitempty"`

	// Pass false to disable unfurling of media content.
	// Optional
	Unfurl_media bool `json:"unfurl_media,omitempty"`

	// URL to an image to use as the icon for this message
	// Optional
	Icon_url string `json:"icon_url,omitempty"`

	// emoji to use as the icon for this message. Overrides icon_url.
	// e.g. :chart_with_upwards_trend:
	// Optional
	Icon_emoji string `json:"icon_emoji,omitempty"`
}

type Attachment struct {

	// text summary of the attachment that is shown by clients that understand
	// attachments but choose not to show them.
	// Required
	Fallback string `json:"fallback"`

	// text that should appear within the attachment
	// Optional
	Text string `json:"text,omitempty"`

	// text that should appear above the formatted data",
	// Optional
	Pretext string `json:"pretext,omitempty"`

	// Can either be one of 'good', 'warning', 'danger', or any hex color code
	// Optional
	Color string `json:"color,omitempty"`

	// Fields are displayed in a table on the message
	Fields []AttachmentFields `json:"fields,omitempty"`
}

func (a *Attachment) String() string {
	str, _ := json.Marshal(a)
	return "[" + string(str) + "]"
}

type AttachmentFields struct {
	// The title may not contain markup and will be escaped for you
	// Required
	Title string `json:"title,omitempty"`

	// Text value of the field. May contain standard message markup and
	// must be escaped as normal. May be multi-line.",
	Value string `json:"value,omitempty"`

	// flag indicating whether the `value` is short enough to be
	// displayed side-by-side with other values
	// Optional
	Short bool `json:"short,omitempty"`
}

func (a *AttachmentFields) String() string {
	str, _ := json.Marshal(a)
	return string(str)
}

// PostMessage sends a message to the channel specified by the id.
//
// Slack API docs: https://api.slack.com/methods/chat.postMessage
func (r *ChannelService) PostMessage(notifReq *SlackMessage) (*http.Response, error) {
	req, err := r.client.NewRequest("POST", BuildQuery(notifReq), nil)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req, nil)
}

func BuildQuery(notifReq *SlackMessage) string {
	var Url *url.URL
	Url, _ = url.Parse(slackEndPoints[base])

	Url.Path += slackEndPoints[chatPostMessage]
	parameters := url.Values{}
	parameters.Add("token", notifReq.Token)
	parameters.Add("channel", notifReq.Channel)
	parameters.Add("text", notifReq.Text)
	parameters.Add("username", notifReq.BotName)
	parameters.Add("parse", notifReq.Parse)
	parameters.Add("link_names", string(notifReq.Link_names))
	parameters.Add("attachment", notifReq.Attachment.String())
	parameters.Add("unfurl_links", "true")
	parameters.Add("unfurl_media", "true")
	parameters.Add("icon_emoji", notifReq.Icon_emoji)
	Url.RawQuery = parameters.Encode()
	return Url.String()
}
