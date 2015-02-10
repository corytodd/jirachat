package jirachat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/corytodd/jirachat"
)

func ExampleNewHipService(r *http.Request) {
	config := HipConfig{
		Token: "<YOUR-ACCESS-TOKEN>",
	}
	svc, err := NewHipService(r, &config)
}

func ExampleNewSlackService(r *http.Request) {
	svc := jirachat.NewSlackService(r,
		&jirachat.SlackConfig{
			ErrChan:    "https://hooks.slack.com/services/<random_webhook_url>",
			BotName:    "Animus",
			WebhookUrl: "https://hooks.slack.com/services/<random_webhook_url>",
			Domain:     "example", // Your JIRA domain, e.g. example.atlassian.net
		})
}

func ExampleResponse() {

	resp := &Response{
		"string": "Hello",
		"int":    42,
	}
	d, _ := json.Marshal(resp)
	fmt.Print(string(d))
	// Output: {"int":42,"string":"Hello"}
}
