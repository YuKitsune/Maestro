package handlers

import (
	"github.com/yukitsune/maestro/pkg/db"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yukitsune/maestro/pkg/api/responses"
	"github.com/yukitsune/maestro/pkg/model"
)

func GetAlbumByIdHandler(repo db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			responses.BadRequest(w, "missing parameter \"id\"")
			return
		}

		foundAlbums, err := repo.GetAlbumsById(r.Context(), id)
		if err != nil {
			responses.Error(w, err)
			return
		}

		if foundAlbums == nil || len(foundAlbums) == 0 {
			responses.NotFoundf(w, "could not find any albums with ID %s", id)
			return
		}

		res := NewResult(model.AlbumType)
		res.AddAll(model.AlbumToHasStreamingServiceSlice(foundAlbums))

		responses.Response(w, res, http.StatusOK)
	}
}
