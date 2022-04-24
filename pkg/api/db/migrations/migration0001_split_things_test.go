package migrations_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/yukitsune/maestro/pkg/api/db/migrations"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func Test_Migration0001ExecutesCorrectly(t *testing.T) {
	withTestDb(t, func(db *mongo.Database) {

		// Seed the database with some data
		err := setupData(db)
		assert.NoError(t, err)

		// Execute the migration
		m := &migrations.Migration0001SplitThings{}
		err = m.Execute(context.Background(), db)
		assert.NoError(t, err)

		// Ensure the database is in the expected state
		err = assertStateIsCorrect(t, db)
		assert.NoError(t, err)
	})
}

func setupData(db *mongo.Database) error {
	thingsColl := db.Collection("things")
	_, err := thingsColl.InsertOne(context.Background(), bson.D{{"name", "my_artist"}, {"thingtype", "artist"}})
	if err != nil {
		return err
	}

	_, err = thingsColl.InsertOne(context.Background(), bson.D{{"name", "my_album"}, {"thingtype", "album"}})
	if err != nil {
		return err
	}

	_, err = thingsColl.InsertOne(context.Background(), bson.D{{"name", "my_track"}, {"thingtype", "track"}})
	if err != nil {
		return err
	}

	return nil
}

func assertStateIsCorrect(t *testing.T, db *mongo.Database) error {

	artistColl := db.Collection("artists")
	c, err := artistColl.CountDocuments(context.Background(), bson.D{{"name", "my_artist"}})
	assert.NoError(t, err)
	assert.Equalf(t, int64(1), c, "artists should be moved to their own collection")

	albumCol := db.Collection("albums")
	c, err = albumCol.CountDocuments(context.Background(), bson.D{{"name", "my_album"}})
	assert.NoError(t, err)
	assert.Equalf(t, int64(1), c, "albums should be moved to their own collection")

	trackCol := db.Collection("tracks")
	c, err = trackCol.CountDocuments(context.Background(), bson.D{{"name", "my_track"}})
	assert.NoError(t, err)
	assert.Equalf(t, int64(1), c, "tracks be moved to their own collection")

	colls, err := db.ListCollectionNames(context.Background(), bson.D{{"options.capped", true}})
	assert.NoError(t, err)

	assert.Falsef(t, contains(colls, "things"), "things collection should be removed")

	return nil
}

func contains(strSlice []string, str string) bool {
	for _, s := range strSlice {
		if s == str {
			return true
		}
	}

	return false
}
