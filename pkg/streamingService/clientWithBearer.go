package streamingService

import (
	"fmt"
	"net/http"
)

type authTransport struct{
	token string
}

func (t *authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	value := fmt.Sprintf("Bearer %s", t.token)
	r.Header.Set("Authorization", value)
	return http.DefaultTransport.RoundTrip(r)
}

func NewClientWithBearer(token string) *http.Client {
	return &http.Client{Transport: &authTransport{token: token}}
}
