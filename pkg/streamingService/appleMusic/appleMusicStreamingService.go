package appleMusic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"maestro/pkg/streamingService"
	"net/http"
	"strings"
)

// https://api.music.apple.com/v1/catalog/us/search

const defaultStorefront = "AU"

type appleMusicStreamingService struct {
	c *http.Client
}

func NewAppleMusicStreamingService(token string) streamingService.StreamingService {
	return &appleMusicStreamingService{c: streamingService.NewClientWithBearerAuth(token)}
}

func (s *appleMusicStreamingService) Name() string {
	return "Apple Music"
}

func (s *appleMusicStreamingService) SearchArtist(name string) (res []streamingService.Artist, err error) {

	term := strings.ReplaceAll(name, " ", "+")
	url := fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/search?term=%s&types=artists", defaultStorefront, term)

	httpRes, err := s.c.Get(url)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return res, err
	}

	var apiRes *SearchResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return res, err
	}

	for _, resource := range apiRes.Results.Artists.Data {
		artist := streamingService.Artist{
			Name:   resource.Attributes.Name,
			Genres: resource.Attributes.GenreNames,
			Url:    resource.Attributes.Url,
		}

		res = append(res, artist)
	}

	return res, nil
}

func (s *appleMusicStreamingService) SearchAlbum(name string) (res []streamingService.Album, err error) {

	term := strings.ReplaceAll(name, " ", "+")
	url := fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/search?term=%s&types=albums", defaultStorefront, term)

	httpRes, err := s.c.Get(url)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return res, err
	}

	var apiRes *SearchResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return res, err
	}

	for _, resource := range apiRes.Results.Albums.Data {
		album := streamingService.Album{
			Name:       resource.Attributes.Name,
			ArtistName: resource.Attributes.ArtistName,
			ArtworkUrl: resource.Attributes.Artwork.Url,
			Url:        resource.Attributes.Url,
		}

		res = append(res, album)
	}

	return res, nil
}

func (s *appleMusicStreamingService) SearchSong(name string) (res []streamingService.Song, err error) {

	term := strings.ReplaceAll(name, " ", "+")
	url := fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/search?term=%s&types=songs", defaultStorefront, term)

	httpRes, err := s.c.Get(url)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return res, err
	}

	var apiRes *SearchResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return res, err
	}

	for _, resource := range apiRes.Results.Songs.Data {
		song := streamingService.Song{
			Name:       resource.Attributes.Name,
			ArtistName: resource.Attributes.ArtistName,
			AlbumName:  resource.Attributes.AlbumName,
			Url:        resource.Attributes.Url,
		}

		res = append(res, song)
	}

	return res, nil
}
