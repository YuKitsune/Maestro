package model

const TrackCollectionName = "tracks"

type Track struct {
	Isrc        string
	Name        string
	ArtistNames []string
	AlbumName   string
	ArtworkLink string

	Source StreamingServiceType
	Market Market
	Link   string
}

func NewTrack(isrc string, name string, artistNames []string, albumName string, artworkLink string, source StreamingServiceType, market Market, link string) *Track {
	return &Track{
		Isrc:        isrc,
		Name:        name,
		ArtistNames: artistNames,
		AlbumName:   albumName,
		ArtworkLink: artworkLink,
		Source:      source,
		Market:      market,
		Link:        link,
	}
}

func (t *Track) GetSource() StreamingServiceType {
	return t.Source
}
