package jirachat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var slackEndPoints = map[API_TARGET]string{
	chatPostMessage: "https://slack.com/api/chat.postMessage",
}

type API_TARGET int

const (
	chatPostMessage API_TARGET = iota
)

type SlackClient struct {
	authToken string
	client    *http.Client
	Room      *ChannelService
}

// NewClient returns a new HipChat API client. You must provide a valid
// AuthToken retrieved from your HipChat account.
func NewClient(authToken string, client *http.Client) *SlackClient {

	c := &SlackClient{
		authToken: authToken,
		client:    client,
	}
	c.Room = &ChannelService{client: c}
	return c
}

// NewRequest creates an API request. This method can be used to performs
// API request not implemented in this library. Otherwise it should not be
// be used directly.
// Relative URLs should always be specified without a preceding slash.
func (c *SlackClient) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	_, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, urlStr, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.authToken)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// Do performs the request, the json received in the response is decoded
// and stored in the value pointed by v.
// Do can be used to perform the request created with NewRequest, as the latter
// it should be used only for API requests not implemented in this library.
func (c *SlackClient) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if c := resp.StatusCode; c < 200 || c > 299 {
		return resp, fmt.Errorf("Server returns status %d", c)
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}
	return resp, err
}
