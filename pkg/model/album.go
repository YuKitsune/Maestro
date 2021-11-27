package model

import "net/url"

const AlbumCollectionKey = "album"

type Album struct {
	Name     string
	ArtistId string

	Services []ServiceSpecificArtist
}

type ServiceSpecificAlbum struct {
	StreamingServiceSpecificEntity
	ArtworkUrl url.URL
}
