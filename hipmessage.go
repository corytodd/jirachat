package jirachat

import (
	"fmt"
	"net/http"
)

// NotificationRequest represents a HipChat room notification request.
type NotificationRequest struct {
	// Background color for message.
	// Valid values: yellow, green, red, purple, gray, random.
	// Defaults to 'yellow'.
	Color string `json:"color,omitempty"`

	// The message body
	// Valid length range: 1 - 10000.
	Message string `json:"message,omitempty"`

	// Whether this message should trigger a user notification (change the tab
	// color, play a sound, notify mobile phones, etc). Each recipient's
	// notification preferences are taken into account.
	Notify bool `json:"notify,omitempty"`

	// Determines how the message is treated by our server and rendered inside
	// HipChat applications
	MessageFormat string `json:"message_format,omitempty"`
}

/// Notification sends a notification to the room specified by the id.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/send_room_notification
func (r *hipService) Notification(id string, notifReq *NotificationRequest) (*http.Response, error) {
	req, err := r.config_.newRequest("POST", fmt.Sprintf("room/%s/notification", id), notifReq)
	if err != nil {
		return nil, err
	}

	return r.config_.do(req, nil)
}
