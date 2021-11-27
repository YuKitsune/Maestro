package model

import "net/url"

const ArtistCollectionKey = "artist"

type Artist struct {
	Name string

	Services []ServiceSpecificArtist
}

type ServiceSpecificArtist struct {
	StreamingServiceSpecificEntity
	ArtworkUrl url.URL
}
