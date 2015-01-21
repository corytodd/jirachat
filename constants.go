package jirachat

import "encoding/json"

const VERSION = "1.3"

// Hipchat specific
const (
	ColorYellow = "yellow"
	ColorRed    = "red"
	ColorGreen  = "green"
	ColorPurple = "purple"
	ColorGray   = "gray"
	ColorRandom = "random"

	FormatText = "text"
	FormatHTML = "html"

	ResponseStatusSent = "sent"
)

// Slack specific
const (
	PARSE_FULL = "full"
)

//Convenience interface for printing anonymous JSON objects
type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}
