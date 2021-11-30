package model

type ThingType string
type ThingGroupId string

const ThingsCollectionName = "things"

const (
	ArtistThing ThingType = "artist"
	AlbumThing ThingType = "album"
	TrackThing ThingType = "track"
)

type Thing interface {
	Type() ThingType
	GetGroupId() ThingGroupId
	SetGroupId(ThingGroupId)
	GetSource() StreamingServiceKey
	GetMarket()Market
 	GetLink() string
}
