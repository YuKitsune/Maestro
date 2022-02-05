package model

type Artist struct {
	Name        string
	ArtworkLink string

	GroupID   ThingGroupID
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

func (a *Artist) GetGroupID() ThingGroupID {
	return a.GroupID
}

func (a *Artist) SetGroupID(groupID ThingGroupID) {
	a.GroupID = groupID
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

func (a *Artist) GetLabel() string {
	return a.Name
}
