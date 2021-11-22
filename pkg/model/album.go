package model

import "net/url"

type Album struct {
	Name string
	ArtistId string

	Services []ServiceSpecificArtist
}

type ServiceSpecificAlbum struct {
	StreamingServiceSpecificEntity
	ArtworkUrl url.URL
}
