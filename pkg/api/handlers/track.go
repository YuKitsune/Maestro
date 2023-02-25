package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/maestro/pkg/db"
	"github.com/yukitsune/maestro/pkg/log"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingservice"
	"net/http"
)

func GetTrackByIsrcHandler(repo db.Repository, serviceProvider streamingservice.ServiceProvider, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLogger, err := log.ForRequest(logger, r)
		if err != nil {
			Error(w, err)
			return
		}

		vars := mux.Vars(r)
		isrc, ok := vars["isrc"]
		if !ok {
			BadRequest(w, "missing parameter \"isrc\"")
			return
		}

		foundTracks, err := repo.GetTracksByIsrc(r.Context(), isrc)
		if err != nil {
			Error(w, err)
			return
		}

		legacyTracks, err := repo.GetTracksByLegacyId(r.Context(), isrc)
		if err != nil {
			Error(w, err)
			return
		}

		for _, legacyTrack := range legacyTracks {
			foundTracks = append(foundTracks, legacyTrack)
		}

		svcs, err := serviceProvider.ListServices()
		if err != nil {
			Error(w, fmt.Errorf("failed to initialize services: %s", err.Error()))
		}

		if len(foundTracks) != len(svcs) {
			newTracks, err := getNewTrackByIsrc(isrc, foundTracks, svcs, reqLogger)
			if err != nil {
				Error(w, err)
				return
			}

			if len(newTracks) > 0 {
				n, err := repo.AddTracks(r.Context(), newTracks)
				if err != nil {
					Error(w, err)
					return
				}

				reqLogger.Infof("%d new tracks added", n)

				for _, newTrack := range newTracks {
					foundTracks = append(foundTracks, newTrack)
				}
			}
		}

		if len(foundTracks) == 0 {
			NotFoundf(w, "could not find any tracks with ISRC code %s", isrc)
			return
		}

		res := NewResult(model.TrackType)
		res.AddAll(model.TrackToHasStreamingServiceSlice(foundTracks))

		Response(w, res, http.StatusOK)
	}
}

func getNewTrackByIsrc(isrc string, knownTracks []*model.Track, svcs streamingservice.StreamingServices, logger *logrus.Entry) ([]*model.Track, error) {

	var tracks []*model.Track

	for key, svc := range svcs {

		// Skip if we know about this track
		trackIsKnown := false
		for _, knownTrack := range knownTracks {
			if knownTrack.Source == key {
				trackIsKnown = true
			}
		}

		if trackIsKnown {
			continue
		}

		track, found, err := svc.GetTrackByIsrc(isrc)
		if err != nil {
			logger.Errorf("%s: %s", key, err.Error())
			continue
		}

		if !found {
			continue
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}
