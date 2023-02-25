package migrations_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/yukitsune/maestro/pkg/db/migrations"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func Test_Migration0002ExecutesCorrectly(t *testing.T) {
	withTestDb(t, func(db *mongo.Database) {

		// Seed the database with some data
		err := setupDataForMigration0002(db)
		assert.NoError(t, err)

		// Execute the migration
		m := &migrations.Migration0002AddThingIds{}
		err = m.Execute(context.Background(), db)
		assert.NoError(t, err)

		// Ensure the database is in the expected state
		err = assertStateIsCorrectForMigration0002(t, db)
		assert.NoError(t, err)
	})
}

func setupDataForMigration0002(db *mongo.Database) error {
	artistsColl := db.Collection("artists")
	_, err := artistsColl.InsertOne(context.Background(), bson.D{{"groupid", "my_artist_id"}})
	if err != nil {
		return err
	}

	albumsColl := db.Collection("albums")
	_, err = albumsColl.InsertOne(context.Background(), bson.D{{"groupid", "my_album_id"}})
	if err != nil {
		return err
	}

	return nil
}

func assertStateIsCorrectForMigration0002(t *testing.T, db *mongo.Database) error {

	artistColl := db.Collection("artists")
	c, err := artistColl.CountDocuments(context.Background(), bson.D{{"artistid", "my_artist_id"}})
	assert.NoError(t, err)
	assert.Equalf(t, int64(1), c, "legacy group IDs should now be artist IDs")

	albumCol := db.Collection("albums")
	c, err = albumCol.CountDocuments(context.Background(), bson.D{{"albumid", "my_album_id"}})
	assert.NoError(t, err)
	assert.Equalf(t, int64(1), c, "legacy group IDs should now be album IDs")

	return nil
}
