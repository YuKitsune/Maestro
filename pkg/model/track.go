package model

type Track struct {
	Name     string
	ArtistId string
	AlbumId  string

	Number int

	Services []ServiceSpecificArtist
}

type ServiceSpecificTrack struct {
	StreamingServiceSpecificEntity
}
