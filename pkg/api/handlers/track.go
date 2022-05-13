package handlers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/maestro/pkg/api/db"
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

		services := serviceProvider.ListServices()
		if len(foundTracks) != len(services) {
			newTracks, err := getNewTrackByIsrc(isrc, foundTracks, services)
			if err != nil {
				Error(w, err)
				return
			}

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

		if len(foundTracks) == 0 {
			NotFoundf(w, "could not find any tracks with ISRC code %s", isrc)
			return
		}

		res := NewResult(model.TrackType)
		res.AddAll(model.TrackToHasStreamingServiceSlice(foundTracks))

		Response(w, res, http.StatusOK)
	}
}

func getNewTrackByIsrc(isrc string, knownTracks []*model.Track, svcs []streamingservice.StreamingService) ([]*model.Track, error) {

	var tracks []*model.Track

	for _, svc := range svcs {

		// Skip if we know about this track
		trackIsKnown := false
		for _, knownTrack := range knownTracks {
			if knownTrack.Source == svc.Key() {
				trackIsKnown = true
			}
		}

		if trackIsKnown {
			continue
		}

		track, found, err := svc.GetTrackByIsrc(isrc)
		if err != nil {
			return tracks, err
		}

		if !found {
			continue
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}
