package jirachat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type EndPoint int

const (
	chatPostMessage EndPoint = iota
	base
	webHook
)

type SlackConfig struct {
	ErrChan    string
	Channel    string
	Token      string
	BotName    string
	Emoji      string
	WebhookUrl string
	client_    http.Client
}

type slackService struct {
	config_ *SlackConfig
}

// Create a new slack with the given config
func NewSlackService(r *http.Request, config *SlackConfig) *slackService {
	client := getHttpClient(r)
	config.client_ = client
	svc := &slackService{config_: config}
	return svc
}

// sendEvent sends Payload which contains JIRA data to Slack.
func (p *Payload) sendEvent(config *SlackConfig) error {
	data, err := json.Marshal(p)
	resp, err := config.client_.Post(config.WebhookUrl, "application/json",
		strings.NewReader(string(data)))
	if err != nil {
		constructSlackError(fmt.Sprintf("%v", err), config.ErrChan)
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
func constructSlackError(msg, channel string) *Payload {
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

	payload := Payload{}
	payload.Channel = channel
	payload.Username = "notify-error"
	payload.Icon_emoji = ":persevere:"
	payload.Unfurl_links = true
	payload.Attachments = []Attachment{attachment}
	payload.Text = ""

	return &payload
}
