package streamingService

import (
	"fmt"
	"net/http"
)

type basicAuthTransport struct{
	token string
}

func (t *basicAuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	value := fmt.Sprintf("Basic %s", t.token)
	r.Header.Set("Authorization", value)
	return http.DefaultTransport.RoundTrip(r)
}

func NewClientWithBasicAuth(token string) *http.Client {
	return &http.Client{Transport: &basicAuthTransport{token: token}}
}
