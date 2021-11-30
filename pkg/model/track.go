package model

type Track struct {
	Name     string
	ArtistNames []string
	AlbumName  string


	GroupId ThingGroupId
	Source StreamingServiceKey
	ThingType ThingType
	Market Market
	Link string
}

func NewTrack(name string, artistNames []string, albumName string, source StreamingServiceKey, market Market, link string) *Track {
	return &Track{
		Name: name,
		ArtistNames: artistNames,
		AlbumName: albumName,
		Source: source,
		ThingType: TrackThing,
		Market: market,
		Link: link,
	}
}

func (t *Track) Type() ThingType {
	return t.ThingType
}

func (t *Track) GetGroupId() ThingGroupId {
	return t.GroupId
}

func (t *Track) SetGroupId(groupId ThingGroupId)  {
	t.GroupId = groupId
}

func (t *Track) GetSource() StreamingServiceKey {
	return t.Source
}

func (t *Track) GetMarket()Market {
	return t.Market
}

func (t *Track) GetLink() string {
	return t.Link
}
