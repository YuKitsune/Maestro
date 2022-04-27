package handlers

import (
	"context"
	"github.com/gorilla/mux"
	mcontext "github.com/yukitsune/maestro/pkg/api/context"
	"github.com/yukitsune/maestro/pkg/api/db"
	"github.com/yukitsune/maestro/pkg/model"
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

	t, err := container.ResolveWithResult(func(ctx context.Context, repo db.Repository) (interface{}, error) {
		foundTracks, err := repo.GetTracksByIsrc(ctx, isrc)
		if err != nil {
			return nil, err
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

	res := &Result{}
	res.AddAll(model.TrackToHasStreamingServiceSlice(tracks))

	Response(w, res, http.StatusOK)
}
