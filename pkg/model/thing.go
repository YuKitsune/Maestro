package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func UnmarshalThingsFromCursor(ctx context.Context, cur *mongo.Cursor) ([]Thing, error) {

	var things []Thing

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		thing, err := UnmarshalThing(cur.Current)
		if err != nil {
			return nil, err
		}

		things = append(things, thing)
	}

	return things, nil
}

func UnmarshalThing(raw bson.Raw) (Thing, error) {

	var thingMap map[string]interface{}
	if err := bson.Unmarshal(raw, &thingMap); err != nil {
		return nil, err
	}

	thingType := ThingType(thingMap["thingtype"].(string))

	switch thingType {
	case ArtistThing:
		var artist *Artist
		if err := bson.Unmarshal(raw, &artist); err != nil{
			return nil, err
		}

		return artist, nil

	case AlbumThing:
		var album *Album
		if err := bson.Unmarshal(raw, &album); err != nil{
			return nil, err
		}

		return album, nil

	case TrackThing:
		var track *Track
		if err := bson.Unmarshal(raw, &track); err != nil{
			return nil, err
		}

		return track, nil

	default:
		return nil, fmt.Errorf("unknown thing type: %s", thingType)
	}
}
