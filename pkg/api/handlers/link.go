package handlers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/camogo"
	mcontext "maestro/pkg/api/context"
	"maestro/pkg/api/db"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"net/http"
	"net/url"
	"reflect"
	"sort"
)

func HandleLink(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	reqLink, ok := vars["link"]
	if !ok {
		BadRequest(w, "missing parameter \"link\"")
		return
	}

	u, err := url.Parse(reqLink)
	if err != nil || u == nil {
		BadRequestf(w, "couldn't parse the given link: %s", reqLink)
		return
	}

	if !u.IsAbs() {
		BadRequest(w, "given link must be absolute")
		return
	}

	container, err := mcontext.Container(r.Context())
	if err != nil {
		Error(w, err)
		return
	}

	// Todo: Add interceptor / middleware support to camogo
	res, err := container.ResolveWithResult(func(logger *logrus.Entry) (interface{}, error) {

		res, err := findForLink(reqLink, container)
		if err != nil {
			logger.Errorf("failed to find things for link %s: %s", reqLink, err.Error())
		}

		return res, err
	})

	if err != nil {
		Error(w, err)
		return
	}

	// Todo: Improve this error message
	if res == nil || len(res.([]model.Thing)) == 0 {
		NotFound(w, "could not find anything")
		return
	}

	Response(w, res, http.StatusOK)
}

func findForLink(link string, container camogo.Container) (interface{}, error) {
	res, err := container.ResolveWithResult(func(ctx context.Context, repo db.Repository, ss []streamingService.StreamingService, logger *logrus.Entry) (interface{}, error) {

		// Trim service-specific stuff from the link
		for _, service := range ss {
			link = service.CleanLink(link)
		}

		logger = logger.WithField("link", link)

		// Search the database for an existing thing with the given link
		foundThing, err := repo.GetThingByLink(ctx, link)
		if err != nil {
			return nil, err
		}

		if foundThing != nil {

			// Link found
			logger = logger.WithField("group_id", foundThing.GetGroupId())
			logger.Debugln("found a thing")

			allThings, err := repo.GetThingsByGroupId(ctx, foundThing.GetGroupId())
			if err != nil {
				return nil, err
			}

			// Check if we're missing any services from our results
			// Todo: It'd be good to have a "Not-found thing" so we can tell if a thing wasn't found for a service,
			// 	rather than assuming it's been newly added
			if len(allThings) < len(ss) {
				logger.Debugf("looks like we have some new services since we found this thing (found %d, looking for %d)\n", len(allThings), len(ss))

				var foundServices []string
				for _, thing := range allThings {
					foundServices = append(foundServices, thing.GetSource().String())
				}

				// Query the remaining streaming service
				sort.Strings(foundServices)
				var newThings []model.Thing
				for _, service := range ss {
					if sort.SearchStrings(foundServices, service.Key().String()) != len(foundServices) {
						continue
					}

					logger.Debugf("fetching thing from %s\n", service.Key())
					thing, err := streamingService.SearchThing(service, foundThing)
					if err != nil {
						return nil, err
					}

					thing.SetGroupId(foundThing.GetGroupId())
					newThings = append(newThings, thing)
				}

				// Add the new things to the database
				if len(newThings) != 0 {
					n, err := repo.AddThings(ctx, newThings)
					if err != nil {
						return nil, err
					}

					logger.Infof("%d new %ss added\n", n, foundThing.Type())
				}
			}

			return allThings, nil
		}

		logger.Debugln("looks like this is a new thing")

		// No links found, query the streaming service and find the same entry on other services
		var targetService streamingService.StreamingService
		var otherServices []streamingService.StreamingService

		// Figure out which streaming service the link belongs to
		err = streamingService.ForEachStreamingService(ss, func(service streamingService.StreamingService) error {
			if service.LinkBelongsToService(link) {
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
			return nil, fmt.Errorf("couldn't find a streaming service for the given link: %s", link)
		}

		// Query the target streaming service
		groupId := model.ThingGroupId(uuid.New().String())
		logger = logger.WithField("group_id", groupId)
		logger.Debugln("using new group id")

		logger.Debugf("searching %s\n", targetService.Key())
		thing, err := targetService.SearchFromLink(link)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", targetService.Key(), err.Error())
		}

		thing.SetGroupId(groupId)

		var things []model.Thing
		things = append(things, thing)

		// Query the other streaming services using what we found from the target streaming service
		err = streamingService.ForEachStreamingService(otherServices, func(service streamingService.StreamingService) error {

			logger.Debugf("searching %s for %s with name %s\n", service.Key(), thing.Type(), thing.GetLabel())
			foundThing, err := streamingService.SearchThing(service, thing)
			if err != nil {
				return fmt.Errorf("%s: %s", service.Key(), err.Error())
			}

			if foundThing == nil || reflect.ValueOf(foundThing).IsNil() {
				return nil
			}

			foundThing.SetGroupId(groupId)
			things = append(things, foundThing)
			return nil
		})
		if err != nil {
			return nil, err
		}

		// Store the new thing in the database
		n, err := repo.AddThings(ctx, things)

		logger.Infof("%d new %ss added\n", n, thing.Type())
		return things, nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}
