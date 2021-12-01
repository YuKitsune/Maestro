package handlers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mcontext "maestro/pkg/api/context"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"net/http"
	"sort"
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

	res, err := container.ResolveWithResult(func(ctx context.Context, db *mongo.Database, ss []streamingService.StreamingService) (interface{}, error) {

		// Trim service-specific stuff from the link
		for _, service := range ss {
			reqLink = service.CleanLink(reqLink)
		}

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
			thingBytes, err := res.DecodeBytes()
			if err != nil {
				return nil, err
			}

			foundThing, err = model.UnmarshalThing(thingBytes)
			if err != nil {
				return nil, err
			}
		}

		if foundThing != nil {
			// Link found

			// Find other things with the same hash
			cur, err := coll.Find(ctx, bson.D{
				{"groupid", foundThing.GetGroupId()},
				{"source", bson.D{{"$ne", foundThing.GetSource()}}},
			})
			if err != nil {
				return nil, err
			}

			var relatedThings []model.Thing
			relatedThings, err = model.UnmarshalThingsFromCursor(ctx, cur)
			if err != nil {
				return nil, err
			}

			allThings := append(relatedThings, foundThing)

			// Check if we're missing any services from our results
			if len(allThings) < len(ss) {
				var foundServices []string
				for _, thing := range allThings {
					foundServices = append(foundServices, thing.GetSource().String())
				}

				// Query the remaining streaming service
				sort.Strings(foundServices)
				var newThings []model.Thing
				for _, service := range ss {
					if sort.SearchStrings(foundServices, service.Name().String()) != len(foundServices) {
						continue
					}

					thing, err := streamingService.SearchThing(service, foundThing)
					if err != nil {
						return nil, err
					}

					thing.SetGroupId(foundThing.GetGroupId())
					newThings = append(newThings, thing)
				}

				// Add the new things to the database
				for _, newThing := range newThings {
					_, err := coll.InsertOne(ctx, newThing)
					if err != nil {
						return nil, err
					}

					allThings = append(allThings, newThing)
				}
			}

			return allThings, nil
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

		if targetService == nil {
			return nil, fmt.Errorf("couldn't find a streaming service for the given link: %s", reqLink)
		}

		// Query the target streaming service
		groupId := model.ThingGroupId(uuid.New().String())
		thing, err := targetService.SearchFromLink(reqLink)
		if err != nil {
			return nil, err
		}

		thing.SetGroupId(groupId)

		// Todo: Strong type would be nice here, but mongo doesn't like it
		var things []interface{}
		things = append(things, thing)

		// Query the other streaming services using what we found from the target streaming service
		err = streamingService.ForEachStreamingService(otherServices, func(service streamingService.StreamingService) error {
			foundThing, err := streamingService.SearchThing(service, thing)
			if err != nil {
				return err
			}

			foundThing.SetGroupId(groupId)
			things = append(things, foundThing)
			return nil
		})
		if err != nil {
			return nil, err
		}

		// Store the new thing in the database
		_, err = coll.InsertMany(ctx, things)
		if err != nil {
			return nil, err
		}

		return things, nil
	})

	if err != nil {
		Error(w, err)
		return
	}

	Response(w, res, http.StatusOK)
}
