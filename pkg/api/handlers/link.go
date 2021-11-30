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

		// Trim service-specific stuff from the link
		for _, service := range ss {
			reqLink  = service.CleanLink(reqLink)
		}

		db := mc.Database(cd.Database)

		var foundThing model.Thing

		// Search the database for an existing thing with the given link
		coll := db.Collection(model.ThingsCollectionName)
		res := coll.FindOne(ctx, bson.D{{"link", reqLink}})
		err = res.Err()
		if err != nil {
			if err != mongo.ErrNoDocuments {
				return nil, res.Err()
			}
		} else {
			err = res.Decode(&foundThing)
			if err != nil {
				return nil, err
			}
		}

		if foundThing != nil {
			// Link found

			// Todo: Find other things with the same hash
			var relatedThings []model.Thing

			// Todo Check if we're missing any services from our results
			if len(relatedThings) < len(ss) {
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

		// Todo: Strong type would be nice here, but mongo doesn't like it
		var things []interface{}
		things = append(things, thing)

		// Query the other streaming services using what we found from the target streaming service
		err = streamingService.ForEachStreamingService(otherServices, func(service streamingService.StreamingService) error {
			foundThing, err := streamingService.SearchThing(service, thing)
			if err != nil {
				return err
			}

			things = append(things, foundThing)
			return nil
		})
		if err != nil {
			return nil, err
		}

		// Store the new thing in the database
		//_, err = coll.InsertMany(ctx, things)
		//if err != nil {
		//	return nil, err
		//}

		return things, nil
	})

	if err != nil {
		Error(w, err)
		return
	}

	Response(w, res, http.StatusOK)
}
