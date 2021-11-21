package deezer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"maestro/pkg/streamingService"
	"net/http"
)

type deezerStreamingService struct {
	client *http.Client
}

func NewDeezerStreamingService() streamingService.StreamingService {
	return &deezerStreamingService{&http.Client{}}
}

func (s *deezerStreamingService) Name() string {
	return "Deezer"
}

func (s *deezerStreamingService) SearchArtist(name string) (res []streamingService.Artist, err error) {

	url := fmt.Sprintf("https://api.deezer.com/search/artist?q=%s", name)

	httpRes, err := s.client.Get(url)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return res, err
	}

	var apiRes *searchArtistResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return res, err
	}

	for _, deezerArtist := range apiRes.Data {
		artist := streamingService.Artist{
			Name: deezerArtist.Name,
			ArtworkUrl: deezerArtist.Picture,
			Url: deezerArtist.Link,
		}

		res = append(res, artist)
	}

	return res, nil
}

func (s *deezerStreamingService) SearchAlbum(name string) (res []streamingService.Album, err error) {

	url := fmt.Sprintf("https://api.deezer.com/search/album?q=%s", name)

	httpRes, err := s.client.Get(url)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return res, err
	}

	var apiRes *searchAlbumResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return res, err
	}

	for _, deezerAlbum := range apiRes.Data {
		album := streamingService.Album{
			Name: deezerAlbum.Title,
			ArtistName: deezerAlbum.Artist.Name,
			ArtworkUrl: deezerAlbum.Cover,
			Url: deezerAlbum.Link,
		}

		res = append(res, album)
	}

	return res, nil
}

func (s *deezerStreamingService) SearchSong(name string) (res []streamingService.Song, err error) {

	url := fmt.Sprintf("https://api.deezer.com/search/track?q=%s", name)

	httpRes, err := s.client.Get(url)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return res, err
	}

	var apiRes *searchTrackResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return res, err
	}

	for _, deezerTrack := range apiRes.Data {
		song := streamingService.Song{
			Name: deezerTrack.Title,
			ArtistName: deezerTrack.Artist.Name,
			AlbumName: deezerTrack.Album.Title,
			Url: deezerTrack.Link,
		}

		res = append(res, song)
	}

	return res, nil
}