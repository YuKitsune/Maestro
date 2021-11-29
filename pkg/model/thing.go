package model

type ThingType string
type ThingHash string

const ThingsCollectionName = "things"

const (
	ArtistThing ThingType = "artist"
	AlbumThing ThingType = "album"
	TrackThing ThingType = "track"
)

type Thing interface {
	Type() ThingType
	GetHash() ThingHash
	GetSource() StreamingServiceKey
	GetMarket()Market
 	GetLink() string
}
