package handlers

import (
	"github.com/gorilla/mux"
	"github.com/yukitsune/maestro/pkg/api/db"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingservice"
	"net/http"
)

// When given a link for a playlist:
// 1. Record the playlist ID and Source, also give it our own ID

// When given a playlist ID
// 1. Fetch the tracks for that playlist

// - We don't need to store tracks, because tracks can change regularly
// - We only need to store the source and ID, so that we can look it up again later
// - Offer an "Import" button for other services

type PlaylistTracks struct {
	Playlist *model.Playlist
	Tracks   []*model.Track
}

func GetPlaylistByIdHandler(repo db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			BadRequest(w, "missing parameter \"id\"")
			return
		}

		foundPlaylist, err := repo.GetPlaylistById(r.Context(), id)
		if err != nil {
			Error(w, err)
			return
		}

		Response(w, foundPlaylist, http.StatusOK)
	}
}

func GetPlaylistTracksByIdHandler(repo db.Repository, serviceProvider streamingservice.ServiceProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			BadRequest(w, "missing parameter \"id\"")
			return
		}

		foundPlaylist, err := repo.GetPlaylistById(r.Context(), id)
		if err != nil {
			Error(w, err)
			return
		}

		svc, err := serviceProvider.GetService(foundPlaylist.GetSource())
		if err != nil {
			Error(w, err)
		}

		tracks, found, err := svc.GetPlaylistTracksById(foundPlaylist.ServicePlaylistId)
		if err != nil {
			Error(w, err)
			return
		}

		if !found {
			NotFoundf(w, "couldn't find tracks for playlist with ID %s in %s", id, foundPlaylist.GetSource())
			return
		}

		pt := &PlaylistTracks{
			foundPlaylist,
			tracks,
		}

		Response(w, pt, http.StatusOK)
	}
}
