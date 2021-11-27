package model

const ArtistCollectionKey = "artist"

type Artist struct {
	Name string

	Links map[StreamingServiceKey]Link
}

func (a *Artist) CollName() string {
	return ArtistCollectionKey
}

func (a *Artist) GetLinks() Links {
	return a.Links
}

func (a *Artist) SetLink(key StreamingServiceKey, link Link) {
	a.Links[key] = link
}
