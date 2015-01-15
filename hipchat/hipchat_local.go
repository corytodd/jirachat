// +build !appengine

package jirachat

import (
	"net/http"
)

//In the local context, return a normal http.Client
func GetHttpClient(r *http.Request) (http.Client, error) {
	var client http.Client
	return client, nil
}
