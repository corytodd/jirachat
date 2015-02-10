// Package hipchat provides a client for using the HipChat API v2.
package jirachat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.hipchat.com/v2/"
)

// Config manages service resources
type HipConfig struct {
	// Hipchat access token
	Token    string
	baseURL_ *url.URL
	client_  *http.Client
}

// HipService gives access to post messages to Hipchat
type hipService struct {
	config_ *HipConfig
}

// NewClient returns a new HipChat API client
func NewHipService(r *http.Request, config *HipConfig) (*hipService, error) {
	baseUrl, err := url.Parse(defaultBaseURL)
	if err != nil {
		panic(err)
	}

	client := getHttpClient(r)
	config.client_ = &client
	config.baseURL_ = baseUrl

	if err = config.IsValid(); err != nil {
		return nil, err
	}

	svc := &hipService{config_: config}
	return svc, err
}

// Returns true if the configuration appears valid
func (c *HipConfig) IsValid() error {
	if len(c.Token) == 0 {
		return errors.New("Invalid Hipchat Token")
	}
	return nil
}

// NewRequest creates an API request. This method can be used to performs
// API request not implemented in this library. Otherwise it should not be
// be used directly.
// Relative URLs should always be specified without a preceding slash.
func (c *HipConfig) newRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.baseURL_.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.Token)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// Do performs the request, the json received in the response is decoded
// and stored in the value pointed by v.
// Do can be used to perform the request created with NewRequest, as the latter
// it should be used only for API requests not implemented in this library.
func (c *HipConfig) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client_.Do(req)
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
