package migrations

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Migration0002AddThingIds struct {
}

func (m *Migration0002AddThingIds) Execute(ctx context.Context, db *mongo.Database) error {

	err := addArtistIds(ctx, db)
	if err != nil {
		return err
	}

	err = addAlbumIds(ctx, db)
	if err != nil {
		return err
	}

	return nil
}

func addArtistIds(ctx context.Context, db *mongo.Database) error {
	artistsColl := db.Collection("artists")

	cur, err := artistsColl.Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var artist bson.M
		if err := bson.Unmarshal(cur.Current, &artist); err != nil {
			return err
		}

		groupId := artist["groupid"].(string)

		_, err := artistsColl.UpdateOne(ctx, bson.D{{"groupid", groupId}}, bson.M{"$set": bson.M{"artistid": groupId}})
		if err != nil {
			return err
		}
	}

	return nil
}

func addAlbumIds(ctx context.Context, db *mongo.Database) error {
	albumsColl := db.Collection("albums")

	cur, err := albumsColl.Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var album bson.M
		if err := bson.Unmarshal(cur.Current, &album); err != nil {
			return err
		}

		groupId := album["groupid"].(string)

		_, err := albumsColl.UpdateOne(ctx, bson.D{{"groupid", groupId}}, bson.M{"$set": bson.M{"albumid": groupId}})
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Migration0002AddThingIds) Version() int {
	return 2
}
