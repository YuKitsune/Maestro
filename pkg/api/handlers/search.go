package handlers

import (
	"github.com/gorilla/mux"
	"github.com/yukitsune/camogo"
	"maestro/pkg/api"
	"maestro/pkg/api/context"
	"maestro/pkg/streamingService"
	"net/http"
)

type SearchArtistResult struct {
	Results map[string][]streamingService.Artist
}

type SearchAlbumResult struct {
	Results map[string][]streamingService.Album
}

type SearchSongResult struct {
	Results map[string][]streamingService.Song
}

func HandleSearchArtist(w http.ResponseWriter, r *http.Request) {

	container, err := context.Container(r.Context())
	if err != nil {
		api.Error(w, err)
		return
	}

	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		api.BadRequest(w, "missing parameter \"name\"")
		return
	}

	res := &SearchArtistResult{}
	res.Results = make(map[string][]streamingService.Artist)

	err = ForEachStreamingService(container, func (service streamingService.StreamingService) error {
		results, err := service.SearchArtist(name)
		if err != nil {
			return err
		}

		res.Results[service.Name()] = results
		return nil
	})
	if err != nil {
		api.Error(w, err)
		return
	}

	api.Response(w, res, http.StatusOK)
}

func HandleSearchAlbum(w http.ResponseWriter, r *http.Request) {

	container, err := context.Container(r.Context())
	if err != nil {
		api.Error(w, err)
		return
	}

	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		api.BadRequest(w, "missing parameter \"name\"")
		return
	}

	res := &SearchAlbumResult{}
	res.Results = make(map[string][]streamingService.Album)

	err = ForEachStreamingService(container, func (service streamingService.StreamingService) error {
		results, err := service.SearchAlbum(name)
		if err != nil {
			return err
		}

		res.Results[service.Name()] = results
		return nil
	})
	if err != nil {
		api.Error(w, err)
		return
	}

	api.Response(w, res, http.StatusOK)
}

func HandleSearchSong(w http.ResponseWriter, r *http.Request) {

	container, err := context.Container(r.Context())
	if err != nil {
		api.Error(w, err)
		return
	}

	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		api.BadRequest(w, "missing parameter \"name\"")
		return
	}

	res := &SearchSongResult{}
	res.Results = make(map[string][]streamingService.Song)

	err = ForEachStreamingService(container, func (service streamingService.StreamingService) error {
		results, err := service.SearchSong(name)
		if err != nil {
			return err
		}

		res.Results[service.Name()] = results
		return nil
	})
	if err != nil {
		api.Error(w, err)
		return
	}

	api.Response(w, res, http.StatusOK)
}

func ForEachStreamingService(container camogo.Container, fn func (streamingService.StreamingService) error) error {
	return container.Resolve(func (services []streamingService.StreamingService) error {
		for _, service := range services {
			err := fn(service)
			if err != nil {
				return err
			}
		}

		return nil
	})
}