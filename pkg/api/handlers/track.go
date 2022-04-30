package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	mcontext "github.com/yukitsune/maestro/pkg/api/context"
	"github.com/yukitsune/maestro/pkg/api/db"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingservice"
	"net/http"
)

func HandleGetTrackByIsrc(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	isrc, ok := vars["isrc"]
	if !ok {
		BadRequest(w, "missing parameter \"isrc\"")
		return
	}

	container, err := mcontext.Container(r.Context())
	if err != nil {
		Error(w, err)
		return
	}

	t, err := container.ResolveWithResult(func(ctx context.Context, repo db.Repository, svcs []streamingservice.StreamingService, logger *logrus.Entry) (interface{}, error) {
		foundTracks, err := repo.GetTracksByIsrc(ctx, isrc)
		if err != nil {
			return nil, err
		}

		legacyTracks, err := repo.GetTracksByLegacyId(ctx, isrc)
		if err != nil {
			return nil, err
		}

		for _, legacyTrack := range legacyTracks {
			foundTracks = append(foundTracks, legacyTrack)
		}

		if len(foundTracks) != len(svcs) {
			newTracks, err := getNewTrackByIsrc(isrc, foundTracks, svcs)
			if err != nil {
				return nil, err
			}

			n, err := repo.AddTracks(ctx, newTracks)
			if err != nil {
				return nil, err
			}

			logger.Infof("%d new tracks added", n)

			for _, newTrack := range newTracks {
				foundTracks = append(foundTracks, newTrack)
			}
		}

		return foundTracks, nil
	})

	if err != nil {
		Error(w, err)
		return
	}

	tracks := t.([]*model.Track)
	if tracks == nil || len(tracks) == 0 {
		NotFoundf(w, "could not find any tracks with ISRC code %s", isrc)
		return
	}

	res := NewResult(model.TrackType)
	res.AddAll(model.TrackToHasStreamingServiceSlice(tracks))

	Response(w, res, http.StatusOK)
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
