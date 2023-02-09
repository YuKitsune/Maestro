package clients

import (
	"fmt"
	"net/http"
)

type bearerAuthTransport struct {
	token string
}

func (t *bearerAuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	value := fmt.Sprintf("Bearer %s", t.token)
	r.Header.Set("Authorization", value)
	return http.DefaultTransport.RoundTrip(r)
}

func NewClientWithBearerAuth(token string) *http.Client {
	return &http.Client{Transport: &bearerAuthTransport{token: token}}
}
