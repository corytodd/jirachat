package jirachat

import (
	"net/http"
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
	BotName string `json:"username"`

	// Change how messages are treated. See below.
	// Optional
	Parse string `json:"parse"`

	// Find and link channel names and usernames.
	// Optional
	Link_names int `json:"link_names"`

	// Structured message attachments
	// Optional
	Fields AttachmentFields `json:"fields"`

	// Optional
	Unfurl_links bool `json:"unfurl_links"`

	// Pass false to disable unfurling of media content.
	// Optional
	Unfurl_media bool `json:"unfurl_media"`

	// URL to an image to use as the icon for this message
	// Optional
	Icon_url string `json:"icon_url"`

	// emoji to use as the icon for this message. Overrides icon_url.
	// e.g. :chart_with_upwards_trend:
	// Optional
	Icon_emoji string `json:"icon_emoji"`
}

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
	Fields AttachmentFields `json:"fields"`
}

type AttachmentFields struct {
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

// PostMessage sends a message to the channel specified by the id.
//
// Slack API docs: https://api.slack.com/methods/chat.postMessage
func (r *ChannelService) PostMessage(id string, notifReq *SlackMessage) (*http.Response, error) {
	req, err := r.client.NewRequest("POST", slackEndPoints[chatPostMessage], notifReq)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req, nil)
}
