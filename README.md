# polychat-jira
JIRA Webhook handler that forwards to Hipchat and/or Slack


This work is mostly taken from https://github.com/tbruyelle/hipchat-go for the Hipchat portions

The Slack integration is adapted from https://github.com/daikikohara/enotify-slack


There are a ton of things that could be done to improve this so feel free to contribute.

If you are using this and I break something, feel free to yell :)

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

	switch event.WebhookEvent {
	case "jira:issue_created":
		err = svc.IssueCreated(event)
	case "jira:issue_deleted":
		err = svc.IssueDeleted(event)
	case "jira:issue_updated":
		err = svc.IssueUpdated(event)
	case "jira:worklog_updated":
		err = svc.WorklogUpdated(event)

	default:
		err = fmt.Errorf("Unknown JIRA Event type %s", event.WebhookEvent)
	}

	if err != nil {
		c.Errorf("Slack Error %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// All went well!
	w.WriteHeader(http.StatusOK)
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
