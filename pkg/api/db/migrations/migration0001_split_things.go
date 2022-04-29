package migrations

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Migration0001SplitThings struct {
}

func (m *Migration0001SplitThings) Execute(ctx context.Context, db *mongo.Database) error {

	thingsColl := db.Collection("things")

	// Migrate artists
	artistCur, err := thingsColl.Find(ctx, bson.D{{"thingtype", "artist"}})
	if err != nil {
		return err
	}

	artists, err := UnmarshalThingsFromCursor(ctx, artistCur)
	if err != nil {
		return err
	}

	if len(artists) > 0 {
		artistsColl := db.Collection("artists")
		_, err = artistsColl.InsertMany(ctx, artists)
		if err != nil {
			return err
		}
	}

	// Migrate albums
	albumsCur, err := thingsColl.Find(ctx, bson.D{{"thingtype", "album"}})
	if err != nil {
		return err
	}

	albums, err := UnmarshalThingsFromCursor(ctx, albumsCur)
	if err != nil {
		return err
	}

	if len(albums) > 0 {
		albumsColl := db.Collection("albums")
		_, err = albumsColl.InsertMany(ctx, albums)
		if err != nil {
			return err
		}
	}

	// Migrate tracks
	tracksCur, err := thingsColl.Find(ctx, bson.D{{"thingtype", "track"}})
	if err != nil {
		return err
	}

	tracks, err := UnmarshalThingsFromCursor(ctx, tracksCur)
	if err != nil {
		return err
	}

	if len(tracks) > 0 {
		tracksColl := db.Collection("tracks")
		_, err = tracksColl.InsertMany(ctx, tracks)
		if err != nil {
			return err
		}
	}

	// Delete things
	//err = thingsColl.Drop(ctx)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (m *Migration0001SplitThings) Version() int {
	return 1
}

func UnmarshalThingsFromCursor(ctx context.Context, cur *mongo.Cursor) ([]interface{}, error) {
	var things []interface{}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var thing interface{}
		if err := bson.Unmarshal(cur.Current, &thing); err != nil {
			return nil, err
		}

		things = append(things, thing)
	}

	return things, nil
}
