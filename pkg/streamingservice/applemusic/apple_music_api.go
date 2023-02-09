package applemusic

import (
	"encoding/json"
	"fmt"
	"github.com/yukitsune/maestro/pkg/clients"
	"io/ioutil"
	"net/http"
	url2 "net/url"

	"github.com/yukitsune/maestro/pkg/model"
)

type ArtistsResult struct {
	Data []*Artist
}

type AlbumResult struct {
	Data []*Album
}

type SongResult struct {
	Data []*Song
}

type Artist struct {
	ID         string            `json:"Id"`
	Attributes *ArtistAttributes //The attributes for the artist.
}

type ArtistAttributes struct {
	GenreNames []string //(Required) The names of the genres associated with this artist.
	Name       string   //(Required) The localized name of the artist.
	URL        string   `json:"Url"` //(Required) The URL for sharing an artist in the iTunes Store.
}

type Song struct {
	ID            string          `json:"Id"`
	Attributes    *SongAttributes //The attributes for the song.
	Relationships *Relationships
}

type SongAttributes struct {
	Isrc        string
	AlbumName   string //(Required) The name of the album the song appears on.
	ArtistName  string //(Required) The artist’s name.
	TrackNumber int    //(Required) The track number.
	Name        string //(Required) The localized name of the song.
	URL         string `json:"Url"` //(Required) The URL for sharing a song in the iTunes Store.
}

type Artwork struct {
	BgColor    string
	Height     int
	Width      int
	TextColor1 string
	TextColor2 string
	TextColor3 string
	TextColor4 string
	URL        string `json:"URl"`
}

type Album struct {
	ID            string           `json:"Id"`
	Attributes    *AlbumAttributes //The attributes for the album.
	Relationships *Relationships
}

type AlbumAttributes struct {
	AlbumName  string  //(Required) The name of the album the music video appears on.
	ArtistName string  //(Required) The artist’s name.
	Artwork    Artwork //The album artwork.
	Name       string  //(Required) The localized name of the album.
	URL        string  `json:"Url"`
	IsSingle   bool
}

type QueryParams struct {
	Term  string
	Types []string
}

type SearchResponse struct {
	Results SearchResult
}

type SearchResult struct {
	Artists *ArtistsResult
	Albums  *AlbumResult
	Songs   *SongResult
}

type Relationships struct {
	Artists struct {
		Data []Relationship
	}
	Albums struct {
		Data []Relationship
	}
}

type Relationship struct {
	Href string
	ID   string `json:"Id"`
	Type string
}

const baseURL = "https://api.music.apple.com"

type client struct {
	client *http.Client
}

func NewAppleMusicClient(token string) *client {
	return &client{client: clients.NewClientWithBearerAuth(token)}
}

func (a *client) SearchArtist(term string, storefront model.Market) ([]Artist, error) {

	querySafeTerm := url2.QueryEscape(term)
	url := fmt.Sprintf("%s/v1/catalog/%s/search?term=%s&types=artists", baseURL, storefront, querySafeTerm)

	httpRes, err := a.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *SearchResponse
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	var artists []Artist
	if res != nil && res.Results.Artists != nil {
		for _, artist := range res.Results.Artists.Data {
			artists = append(artists, *artist)
		}
	}

	return artists, nil
}

func (a *client) SearchAlbum(term string, storefront model.Market) ([]Album, error) {

	querySafeTerm := url2.QueryEscape(term)
	url := fmt.Sprintf("%s/v1/catalog/%s/search?term='%s'&types=albums", baseURL, storefront, querySafeTerm)

	httpRes, err := a.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *SearchResponse
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	var albums []Album
	if res != nil && res.Results.Albums != nil {
		for _, album := range res.Results.Albums.Data {
			albums = append(albums, *album)
		}
	}

	return albums, nil
}

func (a *client) SearchSong(term string, storefront model.Market) ([]Song, error) {

	querySafeTerm := url2.QueryEscape(term)
	url := fmt.Sprintf("%s/v1/catalog/%s/search?term=%s&types=songs", baseURL, storefront, querySafeTerm)

	httpRes, err := a.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *SearchResponse
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	var songs []Song
	if res != nil && res.Results.Songs != nil {
		for _, song := range res.Results.Songs.Data {
			songs = append(songs, *song)
		}
	}

	return songs, nil
}

func (a *client) GetArtist(id string, storefront model.Market) (*Artist, error) {

	url := fmt.Sprintf("%s/v1/catalog/%s/artists/%s", baseURL, storefront, id)

	httpRes, err := a.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *ArtistsResult
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	//var artists []Artist
	//if res != nil && res.Data != nil {
	//	for _, artist := range res.Data {
	//		artists = append(artists, *artist)
	//	}
	//}

	if len(res.Data) == 0 {
		return nil, fmt.Errorf("artist with id %s not found", id)
	}

	artist := res.Data[0]

	return artist, nil
}

func (a *client) GetAlbum(id string, storefront model.Market) (*Album, error) {

	url := fmt.Sprintf("%s/v1/catalog/%s/albums/%s?include=artists", baseURL, storefront, id)

	httpRes, err := a.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *AlbumResult
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	//var albums []Album
	//if res != nil && res.Data != nil {
	//	for _, album := range res.Data {
	//		albums = append(albums, *album)
	//	}
	//}

	if len(res.Data) == 0 {
		return nil, fmt.Errorf("album with id %s not found", id)
	}

	album := res.Data[0]

	return album, nil
}

func (a *client) GetSong(id string, storefront model.Market) (*Song, error) {

	url := fmt.Sprintf("%s/v1/catalog/%s/songs/%s?include=artists,albums", baseURL, storefront, id)

	httpRes, err := a.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *SongResult
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	//var songs []Song
	//if res != nil && res.Data != nil {
	//	for _, song := range res.Data {
	//		songs = append(songs, *song)
	//	}
	//}

	if len(res.Data) == 0 {
		return nil, fmt.Errorf("song with id %s not found", id)
	}

	song := res.Data[0]

	return song, nil
}

func (a *client) GetSongByIsrc(isrc string, storefront model.Market) ([]Song, error) {

	url := fmt.Sprintf("%s/v1/catalog/%s/songs?filter[isrc]=%s", baseURL, storefront, isrc)

	httpRes, err := a.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with %s", httpRes.Status)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)

	var res *SongResult
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	var songs []Song
	for _, song := range res.Data {
		songs = append(songs, *song)
	}

	return songs, nil
}
