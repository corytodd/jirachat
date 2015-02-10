package jirachat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SlackConfig struct {
	// Optional channel to post error reports to
	ErrChan string

	// Receiver channel for JIRA events
	Channel string

	// Bot name reported to Slack
	BotName string

	// Simple Slack Webhook URI
	WebhookUrl string

	// JIRA domain name
	Domain string

	client_ http.Client
}

type slackService struct {
	config_ *SlackConfig
}

// Create a new slack service with the given config. A Slack service
// provides default JIRAWebEvent parser and notification functions.
func NewSlackService(r *http.Request, config *SlackConfig) *slackService {
	client := getHttpClient(r)
	config.client_ = client
	svc := &slackService{config_: config}
	return svc
}

// sendEvent sends Payload which contains JIRA data to Slack.
func (p *payload) sendEvent(config *SlackConfig) error {
	data, err := json.Marshal(p)
	resp, err := config.client_.Post(config.WebhookUrl, "application/json",
		strings.NewReader(string(data)))
	if err != nil {
		constructSlackError(fmt.Sprintf("%v", err), config)
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// ConstructSlackError constructs an error message sent to Slack.
func constructSlackError(msg string, config *SlackConfig) {
	fields := []Field{
		Field{
			Title: "Detail",
			Value: msg,
		},
	}

	attachment := Attachment{
		Fallback: "Error occured on jirachat-slack",
		Pretext:  "Error occured on jirachat-slack",
		Color:    "#FF0000",
		Fields:   fields,
	}

	payload := payload{}
	payload.Username = "Derp Bot"
	payload.Icon_emoji = ":persevere:"
	payload.Unfurl_links = true
	payload.Attachments = []Attachment{attachment}
	payload.Text = ""

	data, _ := json.Marshal(payload)
	config.client_.Post(config.ErrChan, "application/json",
		strings.NewReader(string(data)))
}
