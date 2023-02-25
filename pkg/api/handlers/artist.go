package handlers

import (
	"github.com/yukitsune/maestro/pkg/api/responses"
	"github.com/yukitsune/maestro/pkg/db"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yukitsune/maestro/pkg/model"
)

func GetArtistByIdHandler(repo db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			responses.BadRequest(w, "missing parameter \"id\"")
			return
		}

		foundArtists, err := repo.GetArtistsById(r.Context(), id)
		if err != nil {
			responses.Error(w, err)
			return
		}

		if foundArtists == nil || len(foundArtists) == 0 {
			responses.NotFoundf(w, "could not find any artists with ID %s", id)
			return
		}

		res := NewResult(model.ArtistType)
		res.AddAll(model.ArtistsToHasStreamingServiceSlice(foundArtists))

		responses.Response(w, res, http.StatusOK)
	}
}
