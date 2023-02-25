package model

const ArtistCollectionName = "artists"

type Artist struct {
	ArtistId    string
	Name        string
	ArtworkLink string

	Source StreamingServiceType
	Market Market
	Link   string
}

func NewArtist(name string, artworkLink string, source StreamingServiceType, market Market, link string) *Artist {
	return &Artist{
		Name:        name,
		ArtworkLink: artworkLink,
		Source:      source,
		Market:      market,
		Link:        link,
	}
}

func (a *Artist) GetSource() StreamingServiceType {
	return a.Source
}
