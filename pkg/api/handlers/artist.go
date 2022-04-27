package handlers

import (
	"context"
	"github.com/gorilla/mux"
	mcontext "github.com/yukitsune/maestro/pkg/api/context"
	"github.com/yukitsune/maestro/pkg/api/db"
	"github.com/yukitsune/maestro/pkg/model"
	"net/http"
)

func HandleGetArtistById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		BadRequest(w, "missing parameter \"id\"")
		return
	}

	container, err := mcontext.Container(r.Context())
	if err != nil {
		Error(w, err)
		return
	}

	a, err := container.ResolveWithResult(func(ctx context.Context, repo db.Repository) (interface{}, error) {
		foundArtists, err := repo.GetArtistsById(ctx, id)
		if err != nil {
			return nil, err
		}

		return foundArtists, nil
	})

	if err != nil {
		Error(w, err)
		return
	}

	artists := a.([]*model.Artist)
	if artists == nil || len(artists) == 0 {
		NotFoundf(w, "could not find any artists with ID %s", id)
		return
	}

	res := &Result{}
	res.AddAll(model.ArtistsToHasStreamingServiceSlice(artists))

	Response(w, res, http.StatusOK)
}
