package model

import "net/url"

type Links map[StreamingServiceKey]Link

type Link struct {
	Market    Market
	Url *url.URL
}
