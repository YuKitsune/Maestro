package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mcontext "maestro/pkg/api/context"
	"maestro/pkg/api/db"
	"net/http"
)

func HandleFlag(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	thingType, ok := vars["thingType"]
	thingId, ok := vars["thingId"]
	if !ok {
		BadRequest(w, "missing parameter \"linkId\"")
		return
	}

	container, err := mcontext.Container(r.Context())
	if err != nil {
		Error(w, err)
		return
	}

	err = container.Resolve(func(ctx context.Context, cd *db.Config, mc *mongo.Client) error {
		db := mc.Database(cd.Database)

		// Todo: Don't you dare...
		// 	Validate this...
		coll := db.Collection(thingType)

		_, err := coll.DeleteOne(ctx, bson.D{{"_id", thingId}})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		Error(w, err)
		return
	}

	EmptyResponse(w, http.StatusOK)
}
