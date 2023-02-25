package model

const AlbumCollectionName = "albums"

type Album struct {
	AlbumId     string
	Name        string
	ArtistNames []string
	ArtworkLink string

	Source StreamingServiceType
	Market Market
	Link   string
}

func NewAlbum(name string, artistNames []string, artworkLink string, source StreamingServiceType, market Market, link string) *Album {
	return &Album{
		Name:        name,
		ArtistNames: artistNames,
		ArtworkLink: artworkLink,
		Source:      source,
		Market:      market,
		Link:        link,
	}
}

func (a *Album) GetSource() StreamingServiceType {
	return a.Source
}
