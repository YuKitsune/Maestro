package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mcontext "maestro/pkg/api/context"
	"maestro/pkg/api/errors"
	"maestro/pkg/model"
	"net/http"
)

func HandleGroup(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	groupId, ok := vars["groupId"]
	if !ok {
		BadRequest(w, "missing parameter \"groupId\"")
		return
	}

	container, err := mcontext.Container(r.Context())
	if err != nil {
		Error(w, err)
		return
	}

	res, err := container.ResolveWithResult(func(ctx context.Context, db *mongo.Database) (interface{}, error) {

		var foundThings []model.Thing

		// Search the database for an existing thing with the given link
		coll := db.Collection(model.ThingsCollectionName)
		cur, err := coll.Find(ctx, bson.D{{"groupid", groupId}})

		if err != nil {
			return nil, err
		}

		foundThings, err = model.UnmarshalThingsFromCursor(ctx, cur)
		if err != nil {
			return nil, err
		}

		return foundThings, nil
	})

	if err != nil {
		Error(w, err)
		return
	}

	things := res.([]model.Thing)
	if things == nil || len(things) == 0 {
		Error(w, errors.NotFoundf("could not find group with id %s", groupId))
		return
	}

	Response(w, res, http.StatusOK)
}
