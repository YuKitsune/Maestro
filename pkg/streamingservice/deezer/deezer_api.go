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
	Id      int
	Name    string
	Link    string
	Picture string
}

type searchAlbumResponse struct {
	Data []Album
}

type Album struct {
	Id     int
	Title  string
	Link   string
	Cover  string
	Artist Artist
}

type searchTrackResponse struct {
	Data []Track
}

type Track struct {
	Id     int
	Isrc   string
	Title  string
	Link   string
	Artist Artist
	Album  Album
}

const baseURL = "https://api.deezer.com"

type client struct {
	client *http.Client
}

func NewDeezerClient() *client {
	return &client{client: &http.Client{}}
}

func (d *client) SearchArtist(artistName string) ([]Artist, error) {

	q := url.QueryEscape(fmt.Sprintf("artist:\"%s\"", artistName))
	url := fmt.Sprintf("%s/search/artist?q=%s", baseURL, q)

	httpRes, err := d.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
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

func (d *client) SearchAlbum(artistName string, albumName string) ([]Album, error) {

	q := url.QueryEscape(fmt.Sprintf("artist:\"%s\" album:\"%s\"", artistName, albumName))
	apiURL := fmt.Sprintf("%s/search/album?q=%s", baseURL, q)

	httpRes, err := d.client.Get(apiURL)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
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

func (d *client) SearchTrack(artistName string, albumName string, trackName string) ([]Track, error) {

	q := url.QueryEscape(fmt.Sprintf("artist:\"%s\" album:\"%s\" track:\"%s\"", artistName, albumName, trackName))
	apiURL := fmt.Sprintf("%s/search/track?q=%s", baseURL, q)

	httpRes, err := d.client.Get(apiURL)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
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

func (d *client) GetArtist(id int) (*Artist, error) {

	url := fmt.Sprintf("%s/artist/%d", baseURL, id)

	httpRes, err := d.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *Artist
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *client) GetAlbum(id int) (*Album, error) {

	url := fmt.Sprintf("%s/album/%d", baseURL, id)

	httpRes, err := d.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *Album
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *client) GetTrack(id int) (*Track, error) {

	url := fmt.Sprintf("%s/track/%d", baseURL, id)

	httpRes, err := d.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *Track
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *client) GetTrackByIsrc(isrc string) (*Track, error) {
	url := fmt.Sprintf("%s/track/isrc:%s", baseURL, isrc)

	httpRes, err := d.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *Track
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	// If we haven't got a valid link, don't bother returning it
	if len(res.Link) == 0 {
		return nil, nil
	}

	return res, nil
}
