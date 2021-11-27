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

		var collNames = []string {
			model.TrackCollectionKey,
			model.AlbumCollectionKey,
			model.ArtistCollectionKey,
		}

		// Search the database for an existing thing with the given link
		for _, collName := range collNames {
			coll := db.Collection(collName)

			res := coll.FindOne(ctx, bson.D{{"links.url", reqLink}})
			if res.Err() != nil {
				return nil, res.Err()
			}

			var foundThing model.Thing
			err = res.Decode(&foundThing)
			if err != nil {
				return nil, err
			}

			if foundThing == nil {
				continue
			}

			// Links found

			// Check if we're missing any services from our results

			// Todo: check if we have a streaming service registered that doesn't have a result in our results slice
			if len(foundThing.GetLinks()) < len(ss) {
				// Todo: Query the remaining streaming service
				// Todo: Update the database record for the original thing
			}

			return foundThing, nil
		}

		// No links found, query the streaming service and find the same entry on other services
		var targetService streamingService.StreamingService
		var otherServices []streamingService.StreamingService

		// Figure out which streaming service the link belongs to
		err := streamingService.ForEachStreamingService(ss, func(service streamingService.StreamingService) error {
			if service.LinkBelongsToService(reqLink) {
				targetService = service
			} else {
				otherServices = append(otherServices, service)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		// Query the target streaming service
		thing, err := targetService.SearchFromLink(reqLink)
		if err != nil {
			return nil, err
		}

		thingModel, err := streamingService.ConvertToModel(targetService, thing)
		if err != nil {
			return nil, err
		}

		// Query the other streaming services using what we found from the target streaming service
		err = streamingService.ForEachStreamingService(otherServices, func(service streamingService.StreamingService) error {
			foundThing, err := streamingService.SearchThing(service, thing)
			if err != nil {
				return err
			}

			foundThingModel, err := streamingService.ConvertToModel(service, foundThing)
			if err != nil {
				return err
			}

			// Copy the links
			thingLinks := foundThingModel.GetLinks()
			for key, link := range thingLinks {
				thingModel.SetLink(key, link)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		// Store the new thing in the database
		coll := db.Collection(thingModel.CollName())
		_, err = coll.InsertOne(ctx, thingModel)
		if err != nil {
			return nil, err
		}

		// Todo: Update thing with res.InsertedId? Or does `InsertOne` do that for us?

		return thingModel, nil
	})

	if err != nil {
		Error(w, err)
		return
	}

	Response(w, res, http.StatusOK)
}
