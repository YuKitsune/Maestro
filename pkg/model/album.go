package model

type Album struct {
	Name        string
	ArtistNames []string
	ArtworkLink string

	GroupId   ThingGroupId
	Source    StreamingServiceKey
	ThingType ThingType
	Market    Market
	Link      string
}

func NewAlbum(name string, artistNames []string, artworkLink string, source StreamingServiceKey, market Market, link string) *Album {
	return &Album{
		Name:        name,
		ArtistNames: artistNames,
		ArtworkLink: artworkLink,
		Source:      source,
		ThingType:   AlbumThing,
		Market:      market,
		Link:        link,
	}
}

func (a *Album) Type() ThingType {
	return a.ThingType
}

func (a *Album) GetGroupId() ThingGroupId {
	return a.GroupId
}

func (a *Album) SetGroupId(groupId ThingGroupId) {
	a.GroupId = groupId
}

func (a *Album) GetSource() StreamingServiceKey {
	return a.Source
}

func (a *Album) GetMarket() Market {
	return a.Market
}

func (a *Album) GetLink() string {
	return a.Link
}
