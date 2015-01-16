// Package hipchat provides a client library for the Hipchat REST API.

// +build appengine

package jirachat

import (
	"appengine"
	"appengine/urlfetch"
	"net/http"
)

//In the appengine context, return an http.Client with an appengine context
//using urlfetch transport
func GetHttpClient(r *http.Request) (http.Client, error) {
	return http.Client{Transport: &urlfetch.Transport{Context: appengine.NewContext(r)}}, nil
}
