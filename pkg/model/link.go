package model

import "net/url"

type LinkType string

const (
	LinkTypeArtist LinkType = "artist"
	LinkTypeAlbum = "album"
	LinkTypeTrack = "track"
)

type Link struct {
	Link url.URL
	LinkType LinkType
	EntityId string
}
