package deezer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type searchArtistResponse struct {
	Data []Artist
}

type Artist struct {
	Name    string
	Link    string
	Picture string
}

type searchAlbumResponse struct {
	Data []Album
}

type Album struct {
	Title  string
	Link   string
	Cover  string
	Artist Artist
}

type searchTrackResponse struct {
	Data []Track
}

type Track struct {
	Isrc   string
	Title  string
	Link   string
	Artist Artist
	Album  Album
}

const baseUrl = "https://api.deezer.com"

type DeezerClient struct {
	client *http.Client
}

func NewDeezerClient() *DeezerClient {
	return &DeezerClient{client: &http.Client{}}
}

func (d *DeezerClient) SearchArtist(artistName string) ([]Artist, error) {

	q := url.QueryEscape(fmt.Sprintf("artist:\"%s\"", artistName))
	url := fmt.Sprintf("%s/search/artist?q=%s", baseUrl, q)

	httpRes, err := d.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var apiRes *searchArtistResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return nil, err
	}

	var artists []Artist
	if apiRes != nil && apiRes.Data != nil {
		for _, artist := range apiRes.Data {
			artists = append(artists, artist)
		}
	}

	return artists, nil
}

func (d *DeezerClient) SearchAlbum(artistName string, albumName string) ([]Album, error) {

	q := url.QueryEscape(fmt.Sprintf("artist:\"%s\" album:\"%s\"", artistName, albumName))
	apiUrl := fmt.Sprintf("%s/search/album?q=%s", baseUrl, q)

	httpRes, err := d.client.Get(apiUrl)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var apiRes *searchAlbumResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return nil, err
	}

	var albums []Album
	if apiRes != nil && apiRes.Data != nil {
		for _, album := range apiRes.Data {
			albums = append(albums, album)
		}
	}

	return albums, nil
}

func (d *DeezerClient) SearchTrack(artistName string, albumName string, trackName string) ([]Track, error) {

	q := url.QueryEscape(fmt.Sprintf("artist:\"%s\" album:\"%s\" track:\"%s\"", artistName, albumName, trackName))
	apiUrl := fmt.Sprintf("%s/search/track?q=%s", baseUrl, q)

	httpRes, err := d.client.Get(apiUrl)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var apiRes *searchTrackResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return nil, err
	}

	var tracks []Track
	if apiRes != nil && apiRes.Data != nil {
		for _, track := range apiRes.Data {
			tracks = append(tracks, track)
		}
	}

	return tracks, nil
}

func (d *DeezerClient) GetArtist(id string) (*Artist, error) {

	url := fmt.Sprintf("%s/artist/%s", baseUrl, id)

	httpRes, err := d.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *Artist
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *DeezerClient) GetAlbum(id string) (*Album, error) {

	url := fmt.Sprintf("%s/album/%s", baseUrl, id)

	httpRes, err := d.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *Album
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *DeezerClient) GetTrack(id string) (*Track, error) {

	url := fmt.Sprintf("%s/track/%s", baseUrl, id)

	httpRes, err := d.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *Track
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *DeezerClient) GetTrackByIsrc(isrc string) (*Track, error) {
	url := fmt.Sprintf("%s/track/isrc:%s", baseUrl, isrc)

	httpRes, err := d.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *Track
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
