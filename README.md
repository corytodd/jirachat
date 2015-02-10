# polychat-jira
JIRA Webhook handler that forwards to Hipchat and/or Slack


This work is mostly taken from https://github.com/tbruyelle/hipchat-go for the Hipchat portions

The Slack integration is adapted from https://github.com/daikikohara/enotify-slack


There are a ton of things that could be done to improve this so feel free to contribute.

If you are using this and I break something, feel free to yell :). 

How to work with the Slack service
```
// Slack Sample Handler
func SendSlack(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	svc := jirachat.NewSlackService(r,
		&jirachat.SlackConfig{
			ErrChan:    <YOU_ERROR_CHANNEL>,
			BotName:    <BOT_NAME>,
			WebhookUrl: <SLACK_WEBHOOK_URL>,
		})

	// Parse our event, baby! JIRA event can be touchy,
	// don't consider parse errors fatal
	event, err := jirachat.Parse(r)
	if err != nil {
		c.Errorf("Error parsing JIRA event %v", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//return
	}

	mySvc := &MySlacker{config: svc.Config}
	err = mySvc.IssueCreated(&event)

	if err != nil {
		c.Errorf("Slack Error %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// All went well!
	w.WriteHeader(http.StatusOK)
}

type MySlacker struct {
	config *jirachat.SlackConfig
}

func (s *MySlacker) IssueCreated(event *jirachat.JIRAWebevent) error {
	payload := jirachat.SlackMessage{}
	fields := []jirachat.Field{
		jirachat.Field{
			Title: "Summary",
			Value: event.Issue.Fields.Summary,
			Short: false,
		},
	}
	title := fmt.Sprintf("%s created %s", event.GetUserLink(s.config),
		event.GetIssueLink(s.config))
	attachment := jirachat.Attachment{
		Fallback: title,
		Pretext:  title,
		Color:    getPriorityColor(event.Issue.Fields.Priority.Id),
		Fields:   fields,
	}

	payload.Channel = s.config.Channel
	payload.Username = s.config.BotName
	payload.Icon_url = event.User.LargeAvatar()
	payload.Unfurl_links = true
	payload.Text = ""
	payload.Attachments = []jirachat.Attachment{attachment}
	return payload.SendEvent(s.config)
}
func getPriorityColor(id string) string {

	switch {
	case id == "1": 		// Blocker
		return "#990000"
	case id == "2":
		return "#cc0000" 	// Critical
	case id == "3":
		return "#ff0000"
	case id == "6": 		// Normal
		return "#339933"
	case id == "4": 		// Minor
		return "#006600"
	case id == "5": 		// Trivial
		return "#003300"
	case id == "10000": 	// Holding
		return "#000000"
	default:
		return jriachat.ColorGood
	}
}

```

How to work with the Hipchat service
```
// Sample Hipchat Handler
func SendHipchat(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	client, err := jirachat.NewHipService(r,
		&jirachat.HipConfig{Token: <YOUR_API_TOKEN>})

	if err != nil {
		c.Errorf("Failed to create HipService: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req jirachat.NotificationRequest
	var sendIt bool

	// Parse our event, baby! JIRA event can be touchy,
	// don't consider parse errors fatal
	event, err := jirachat.Parse(r)
	if err != nil {
		c.Errorf("Error parsing JIRA event %v", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//return
	}

	switch event.WebhookEvent {
	case "jira:issue_created":
		req = jirachat.NotificationRequest{
			Message: getUserAvatar(event.User) + " " +
				getUserLink(event.User) +
				"<i><strong> created </i> &#8594</strong> " +
				getLink(event.Issue),
			Color:         jirachat.ColorGreen,
			MessageFormat: jirachat.FormatHTML,
			Notify:        true,
		}
	}
	sendIt = true
	if sendIt {
		if resp, err := client.Notification(<ROOM_ID>, &req); err != nil {
			c.Errorf(err.Error())
			fmt.Printf("Error during room notification %q\n", err)
			fmt.Printf("Server returns %+v\n", resp)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Everything went well!
		w.WriteHeader(http.StatusOK)
	}
}
```
