package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mcontext "maestro/pkg/api/context"
	"maestro/pkg/model"
	"net/http"
)

// Todo: Need to refine the flagging process

func HandleFlag(w http.ResponseWriter, r *http.Request) {

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

	err = container.Resolve(func(ctx context.Context, db *mongo.Database) error {

		coll := db.Collection(model.ThingsCollectionName)

		_, err := coll.DeleteMany(ctx, bson.D{{"groupid", groupId}})
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
