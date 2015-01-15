package jirachat

// RoomService gives access to the room related methods of the API.
type SlackService struct {
	client *Client
}

type Message struct {
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

	// Pass true to enable unfurling of primarily text-based content.
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
