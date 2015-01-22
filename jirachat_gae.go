// Package hipchat provides a client library for the Hipchat REST API.

// +build appengine

package jirachat

import (
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

//In the appengine context, return an http.Client with an appengine context
//using urlfetch transport
func getHttpClient(r *http.Request) http.Client {
	return http.Client{Transport: &urlfetch.Transport{Context: appengine.NewContext(r)}}
}
