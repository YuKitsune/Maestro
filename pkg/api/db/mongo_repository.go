package db

import (
	"context"
	"github.com/yukitsune/maestro/pkg/metrics"
	"github.com/yukitsune/maestro/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	db  *mongo.Database
	rec metrics.Recorder
}

func NewMongoRepository(db *mongo.Database, rec metrics.Recorder) Repository {
	return &mongoRepository{db, rec}
}

func (m mongoRepository) GetThingByLink(ctx context.Context, s string) (model.Thing, error) {
	go m.rec.CountDatabaseCall()

	var foundThing model.Thing
	coll := m.db.Collection(model.ThingsCollectionName)
	res := coll.FindOne(ctx, bson.D{{"link", s}})
	err := res.Err()
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	} else {
		thingBytes, err := res.DecodeBytes()
		if err != nil {
			return nil, err
		}

		foundThing, err = model.UnmarshalThing(thingBytes)
		if err != nil {
			return nil, err
		}
	}

	return foundThing, nil
}

func (m mongoRepository) GetThingsByGroupId(ctx context.Context, id model.ThingGroupId) ([]model.Thing, error) {
	go m.rec.CountDatabaseCall()

	coll := m.db.Collection(model.ThingsCollectionName)
	cur, err := coll.Find(ctx, bson.D{
		{"groupid", id},
		// {"source", bson.D{{"$ne", foundThing.GetSource()}}},
	})

	var things []model.Thing
	things, err = model.UnmarshalThingsFromCursor(ctx, cur)
	if err != nil {
		return nil, err
	}

	return things, nil
}

func (m mongoRepository) AddThings(ctx context.Context, things []model.Thing) (int, error) {
	go m.rec.CountDatabaseCall()

	coll := m.db.Collection(model.ThingsCollectionName)
	insertRes, err := coll.InsertMany(ctx, thingsToInterfaces(things))
	if err != nil {
		return 0, err
	}

	return len(insertRes.InsertedIDs), nil
}

func thingsToInterfaces(things []model.Thing) []interface{} {
	var s []interface{}
	for _, thing := range things {
		s = append(s, thing)
	}

	return s
}
