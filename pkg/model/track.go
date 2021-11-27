package model

const TrackCollectionKey = "track"

type Track struct {
	Name     string
	ArtistId string
	AlbumId  string

	Number int

	Links map[StreamingServiceKey]Link
}

func (t *Track) CollName() string {
	return TrackCollectionKey
}

func (t *Track) GetLinks() Links {
	return t.Links
}

func (t *Track) SetLink(key StreamingServiceKey, link Link) {
	t.Links[key] = link
}
