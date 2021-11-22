package model

import "net/url"

type StreamingServiceKey string

type StreamingServiceSpecificEntity struct {
	ServiceKey StreamingServiceKey
	Country Country

	Id string
	Url url.URL
}