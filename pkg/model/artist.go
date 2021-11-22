package model

import "net/url"

type Artist struct {
	Name string

	Services []ServiceSpecificArtist
}

type ServiceSpecificArtist struct {
	StreamingServiceSpecificEntity
	ArtworkUrl url.URL
}
