// +build !appengine

package jirachat

import (
	"net/http"
)

//In the local context, return a normal http.Client
func getHttpClient(r *http.Request) http.Client {
	var client http.Client
	return client
}
