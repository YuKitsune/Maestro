package handlers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/camogo"
	mcontext "github.com/yukitsune/maestro/pkg/api/context"
	"github.com/yukitsune/maestro/pkg/api/db"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingservice"
	"net/http"
	"net/url"
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

	if res == nil || len(res.([]model.Thing)) == 0 {
		NotFound(w, "could not find anything")
		return
	}

	Response(w, res, http.StatusOK)
}

func findForLink(link string, container camogo.Container) (interface{}, error) {
	res, err := container.ResolveWithResult(func(ctx context.Context, repo db.Repository, ss []streamingservice.StreamingService, logger *logrus.Entry) (interface{}, error) {

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
			things, err := findForExistingThing(foundThing, container)
			return things, err
		}

		logger.Debugln("looks like this is a new thing")

		// No links found, query the streaming service and find the same entry on other services
		var targetService streamingservice.StreamingService
		var otherServices []streamingservice.StreamingService

		// Figure out which streaming service the link belongs to
		for _, service := range ss {
			if service.LinkBelongsToService(link) {
				targetService = service
			} else {
				otherServices = append(otherServices, service)
			}
		}

		if targetService == nil {
			return nil, fmt.Errorf("couldn't find a streaming service for the given link: %s", link)
		}

		// Since this is a new thing, we need a new ID
		groupID := model.ThingGroupID(uuid.New().String())
		logger = logger.WithField("group_id", groupID)
		logger.Debugln("using new group id")

		// Query the target streaming service
		logger.Debugf("searching %s\n", targetService.Key())
		thing, found, err := targetService.SearchFromLink(link)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", targetService.Key(), err.Error())
		}

		if !found {
			return nil, fmt.Errorf("%s: couldn't find anything for the given link", targetService.Key())
		}

		thing.SetGroupID(groupID)

		var things []model.Thing
		things = append(things, thing)

		// Query the other streaming services using what we found from the target streaming service
		for _, service := range otherServices {
			logger.Debugf("searching %s for %s with name %s\n", service.Key(), thing.Type(), thing.GetLabel())
			foundThing, found, err := streamingservice.SearchThing(service, thing)
			if err != nil {
				return nil, fmt.Errorf("%s: %s", service.Key(), err.Error())
			}

			if !found {
				logger.Debugf("couldn't find anything for %s", service.Key())
				continue
			}

			foundThing.SetGroupID(groupID)
			things = append(things, foundThing)
		}

		// Store the new thing in the database
		n, err := repo.AddThings(ctx, things)

		logger.Infof("%d new %ss added\n", n, thing.Type())
		return things, nil
	})

	return res, err
}

func findForExistingThing(foundThing model.Thing, container camogo.Container) (interface{}, error) {
	return container.ResolveWithResult(func(ctx context.Context, repo db.Repository, ss []streamingservice.StreamingService, logger *logrus.Entry) (interface{}, error) {

		logger = logger.WithField("group_id", foundThing.GetGroupID())
		logger.Debugln("found a thing")

		// Find any related things based on the group ID
		// Todo: If we've found a track, then query the DB with the ISRC code
		// 	Otherwise, use the group ID
		existingThings, err := repo.GetThingsByGroupID(ctx, foundThing.GetGroupID())
		if err != nil {
			return nil, err
		}

		// Check if we're missing any services from our results
		// Todo: It'd be good to have a "Not-found thing" so we can tell if a thing wasn't found for a service,
		// 	rather than assuming it's been newly added
		if len(existingThings) >= len(ss) {
			return existingThings, nil
		}
		logger.Debugf("looks like we have some new services since we found this thing (found %d, looking for %d)\n", len(existingThings), len(ss))

		// Query the remaining streaming service
		var newThings []model.Thing
		for _, service := range ss {
			if serviceIsInThings(existingThings, service.Key()) {
				continue
			}

			logger.Debugf("fetching thing from %s\n", service.Key())
			thing, found, err := streamingservice.SearchThing(service, foundThing)
			if err != nil {
				return nil, err
			}

			if !found {
				logger.Debugf("couldn't find anything for %s", service.Key())
				continue
			}

			thing.SetGroupID(foundThing.GetGroupID())
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

		allThings := append(existingThings, newThings...)
		return allThings, nil
	})
}

func serviceIsInThings(things []model.Thing, service model.StreamingServiceKey) bool {
	for _, thing := range things {
		if thing.GetSource() == service {
			return true
		}
	}

	return false
}
