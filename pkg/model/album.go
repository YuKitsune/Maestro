package model

const AlbumCollectionKey = "album"

type Album struct {
	Name     string
	ArtistId string

	Links Links
}

func (a *Album) CollName() string {
	return AlbumCollectionKey
}

func (a *Album) GetLinks() Links {
	return a.Links
}

func (a *Album) SetLink(key StreamingServiceKey, link Link) {
	a.Links[key] = link
}
