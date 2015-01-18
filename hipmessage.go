package jirachat

import (
	"fmt"
	"net/http"
)

// NotificationRequest represents a HipChat room notification request.
type NotificationRequest struct {
	Color         string `json:"color,omitempty"`
	Message       string `json:"message,omitempty"`
	Notify        bool   `json:"notify,omitempty"`
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
