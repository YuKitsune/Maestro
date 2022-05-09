package handlers

import (
	"context"
	"github.com/gorilla/mux"
	mcontext "github.com/yukitsune/maestro/pkg/api/context"
	"github.com/yukitsune/maestro/pkg/api/db"
	"github.com/yukitsune/maestro/pkg/model"
	"net/http"
)

func HandleGetAlbumById(w http.ResponseWriter, r *http.Request) {

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
		foundAlbums, err := repo.GetAlbumsById(ctx, id)
		if err != nil {
			return nil, err
		}

		return foundAlbums, nil
	})

	if err != nil {
		Error(w, err)
		return
	}

	albums := a.([]*model.Album)
	if albums == nil || len(albums) == 0 {
		NotFoundf(w, "could not find any albums with ID %s", id)
		return
	}

	res := NewResult(model.AlbumType)
	res.AddAll(model.AlbumToHasStreamingServiceSlice(albums))

	Response(w, res, http.StatusOK)
}
