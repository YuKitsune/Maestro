package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mcontext "maestro/pkg/api/context"
	"maestro/pkg/api/db"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"net/http"
)

func HandleLink(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	reqLink, ok := vars["link"]
	if !ok {
		BadRequest(w, "missing parameter \"link\"")
		return
	}

	container, err := mcontext.Container(r.Context())
	if err != nil {
		Error(w, err)
		return
	}

	res, err := container.ResolveWithResult(func(ctx context.Context, cd *db.Config, mc *mongo.Client, ss []streamingService.StreamingService) (interface{}, error) {
		db := mc.Database(cd.Database)
		coll := db.Collection("links")

		// Todo: Strong type would be sick
		var results []streamingService.Thing

		var foundLinks []model.Link
		cur, err := coll.Find(ctx, bson.D{{"link", reqLink}})
		if err != nil {
			return nil, err
		}

		err = cur.All(ctx, &foundLinks)
		if err != nil {
			return nil, err
		}

		if len(foundLinks) == 0 {
			// No links found, query the streaming service and find the same entry on other services

			var targetService streamingService.StreamingService
			var otherServices []streamingService.StreamingService

			err := streamingService.ForEachStreamingService(ss, func(service streamingService.StreamingService) error {
				if service.LinkBelongsToService(reqLink) {
					targetService = service
				} else {
					otherServices = append(otherServices, service)
				}
				return nil
			})

			if err != nil {
				return results, err
			}

			var foundThings []streamingService.Thing
			thing, err := targetService.SearchFromLink(reqLink)
			if err != nil {
				return results, err
			}

			foundThings = append(foundThings, thing)

			err = streamingService.ForEachStreamingService(otherServices, func(service streamingService.StreamingService) error {
				foundThing, err := streamingService.SearchThing(service, thing)
				if err != nil {
					return err
				}

				foundThings = append(foundThings, foundThing)
				return nil
			})
			if err != nil {
				return results, err
			}

			// Todo: Store the results in the database

			results = foundThings
			return results, nil

		} else {
			// Links found,
			// Todo: get the full data from the relative table

			// Check if we're missing any services from our results

			// Todo: check if we have a streaming service registered that doesn't have a result in our results slice
			if len(foundLinks) < len(ss) {
				// Todo: Query the remaining streaming service
				// Todo: Store the results in the database
				// Todo: Add the results to the results slice
			}
		}

		return results, nil
	})
	if err != nil {
		Error(w, err)
		return
	}

	Response(w, res, http.StatusOK)
}

func HandleFlagLink(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	linkId, ok := vars["linkId"]
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
		coll := db.Collection("links")

		// Todo: Consider deleting related things too

		_, err := coll.DeleteOne(ctx, bson.D{{"_id", linkId}})
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
