package handlers

import (
	"github.com/gorilla/mux"
	"maestro/pkg/api/context"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"net/http"
)

type SearchArtistResult struct {
	Results map[model.StreamingServiceKey]*model.Artist
}

type SearchAlbumResult struct {
	Results map[model.StreamingServiceKey]*model.Album
}

type SearchSongResult struct {
	Results map[model.StreamingServiceKey]*model.Track
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
	res.Results = make(map[model.StreamingServiceKey]*model.Artist)

	err = container.Resolve(func(services []streamingService.StreamingService) error {
		return streamingService.ForEachStreamingService(services, func(service streamingService.StreamingService) error {
			results, err := service.SearchArtist(&model.Artist{
				Name:   name,
				Market: model.DefaultMarket,
			})
			if err != nil {
				return err
			}

			res.Results[service.Key()] = results
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
	res.Results = make(map[model.StreamingServiceKey]*model.Album)

	err = container.Resolve(func(services []streamingService.StreamingService) error {
		return streamingService.ForEachStreamingService(services, func(service streamingService.StreamingService) error {
			results, err := service.SearchAlbum(&model.Album{
				Name:   name,
				Market: model.DefaultMarket,
			})
			if err != nil {
				return err
			}

			res.Results[service.Key()] = results
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
	res.Results = make(map[model.StreamingServiceKey]*model.Track)

	err = container.Resolve(func(services []streamingService.StreamingService) error {
		return streamingService.ForEachStreamingService(services, func(service streamingService.StreamingService) error {
			results, err := service.SearchSong(&model.Track{
				Name:   name,
				Market: model.DefaultMarket,
			})
			if err != nil {
				return err
			}

			res.Results[service.Key()] = results
			return nil
		})
	})
	if err != nil {
		Error(w, err)
		return
	}

	Response(w, res, http.StatusOK)
}
