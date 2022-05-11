package handlers

import (
	"github.com/gorilla/mux"
	"github.com/yukitsune/maestro/pkg/api/db"
	"github.com/yukitsune/maestro/pkg/model"
	"net/http"
)

func GetArtistByIdHandler(repo db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			BadRequest(w, "missing parameter \"id\"")
			return
		}

		foundArtists, err := repo.GetArtistsById(r.Context(), id)
		if err != nil {
			Error(w, err)
			return
		}

		if foundArtists == nil || len(foundArtists) == 0 {
			NotFoundf(w, "could not find any artists with ID %s", id)
			return
		}

		res := NewResult(model.ArtistType)
		res.AddAll(model.ArtistsToHasStreamingServiceSlice(foundArtists))

		Response(w, res, http.StatusOK)
	}
}
