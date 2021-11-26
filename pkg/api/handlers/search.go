package handlers

import (
	"github.com/gorilla/mux"
	"maestro/pkg/api/context"
	"maestro/pkg/streamingService"
	"net/http"
)

const defaultRegion streamingService.Region = "AU"

type SearchArtistResult struct {
	Results map[string]*streamingService.Artist
}

type SearchAlbumResult struct {
	Results map[string]*streamingService.Album
}

type SearchSongResult struct {
	Results map[string]*streamingService.Song
}

func HandleSearchArtist(w http.ResponseWriter, r *http.Request) {

	container, err := context.Container(r.Context())
	if err != nil {
		Error(w, err)
		return
	}

	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		BadRequest(w, "missing parameter \"name\"")
		return
	}

	res := &SearchArtistResult{}
	res.Results = make(map[string]*streamingService.Artist)

	err = container.Resolve(func (services []streamingService.StreamingService) error {
		return streamingService.ForEachStreamingService(services, func(service streamingService.StreamingService) error {
			results, err := service.SearchArtist(&streamingService.Artist{
				Name: name,
				Region: defaultRegion,
			})
			if err != nil {
				return err
			}

			res.Results[service.Name()] = results
			return nil
		})
	})

	if err != nil {
		Error(w, err)
		return
	}

	Response(w, res, http.StatusOK)
}

func HandleSearchAlbum(w http.ResponseWriter, r *http.Request) {

	container, err := context.Container(r.Context())
	if err != nil {
		Error(w, err)
		return
	}

	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		BadRequest(w, "missing parameter \"name\"")
		return
	}

	res := &SearchAlbumResult{}
	res.Results = make(map[string]*streamingService.Album)

	err = container.Resolve(func (services []streamingService.StreamingService) error {
		return streamingService.ForEachStreamingService(services, func(service streamingService.StreamingService) error {
			results, err := service.SearchAlbum(&streamingService.Album{
				Name: name,
				Region: defaultRegion,
			})
			if err != nil {
				return err
			}

			res.Results[service.Name()] = results
			return nil
		})
	})
	if err != nil {
		Error(w, err)
		return
	}

	Response(w, res, http.StatusOK)
}

func HandleSearchSong(w http.ResponseWriter, r *http.Request) {

	container, err := context.Container(r.Context())
	if err != nil {
		Error(w, err)
		return
	}

	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		BadRequest(w, "missing parameter \"name\"")
		return
	}

	res := &SearchSongResult{}
	res.Results = make(map[string]*streamingService.Song)

	err = container.Resolve(func (services []streamingService.StreamingService) error {
		return streamingService.ForEachStreamingService(services, func(service streamingService.StreamingService) error {
			results, err := service.SearchSong(&streamingService.Song{
				Name: name,
				Region: defaultRegion,
			})
			if err != nil {
				return err
			}

			res.Results[service.Name()] = results
			return nil
		})
	})
	if err != nil {
		Error(w, err)
		return
	}

	Response(w, res, http.StatusOK)
}
