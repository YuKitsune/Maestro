package model

type Artist struct {
	Name        string
	ArtworkLink string

	GroupId   ThingGroupId
	Source    StreamingServiceKey
	ThingType ThingType
	Market    Market
	Link      string
}

func NewArtist(name string, artworkLink string, source StreamingServiceKey, market Market, link string) *Artist {
	return &Artist{
		Name:        name,
		ArtworkLink: artworkLink,
		Source:      source,
		ThingType:   ArtistThing,
		Market:      market,
		Link:        link,
	}
}

func (a *Artist) Type() ThingType {
	return a.ThingType
}

func (a *Artist) GetArtworkLink() string {
	return a.ArtworkLink
}

func (a *Artist) GetGroupId() ThingGroupId {
	return a.GroupId
}

func (a *Artist) SetGroupId(groupId ThingGroupId) {
	a.GroupId = groupId
}

func (a *Artist) GetSource() StreamingServiceKey {
	return a.Source
}

func (a *Artist) GetMarket() Market {
	return a.Market
}

func (a *Artist) GetLink() string {
	return a.Link
}
