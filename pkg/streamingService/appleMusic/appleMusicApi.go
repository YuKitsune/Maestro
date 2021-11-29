package appleMusic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"net/http"
	"strings"
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
	Id         string
	Attributes *ArtistAttributes //The attributes for the artist.
}

type ArtistAttributes struct {
	GenreNames []string //(Required) The names of the genres associated with this artist.
	Name       string   //(Required) The localized name of the artist.
	Url        string   //(Required) The URL for sharing an artist in the iTunes Store.
}

type Song struct {
	Id         string
	Attributes *SongAttributes //The attributes for the song.
	Relationships *Relationships
}

type SongAttributes struct {
	AlbumName  string //(Required) The name of the album the song appears on.
	ArtistName string //(Required) The artist’s name.
	TrackNumber int //(Required) The track number.
	Name       string //(Required) The localized name of the song.
	Url        string //(Required) The URL for sharing a song in the iTunes Store.
}

type Artwork struct {
	BgColor    string
	Height     int
	Width      int
	TextColor1 string
	TextColor2 string
	TextColor3 string
	TextColor4 string
	Url        string
}

type Album struct {
	Id         string
	Attributes *AlbumAttributes //The attributes for the album.
	Relationships *Relationships
}

type AlbumAttributes struct {
	AlbumName  string  //(Required) The name of the album the music video appears on.
	ArtistName string  //(Required) The artist’s name.
	Artwork    Artwork //The album artwork.
	Name       string  //(Required) The localized name of the album.
	Url        string
	IsSingle bool
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
	Id string
	Type string
}

const baseUrl = "https://api.music.apple.com"

type AppleMusicClient struct {
	client *http.Client
}

func NewAppleMusicClient(token string) *AppleMusicClient {
	return &AppleMusicClient{client: streamingService.NewClientWithBearerAuth(token)}
}

func (a *AppleMusicClient) SearchArtist(term string, storefront model.Market) ([]Artist, error) {

	querySafeTerm := strings.ReplaceAll(term, " ", "+")
	url := fmt.Sprintf("%s/v1/catalog/%s/search?term=%s&types=artists", baseUrl, storefront, querySafeTerm)

	httpRes, err := a.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
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

func (a *AppleMusicClient) SearchAlbum(term string, storefront model.Market) ([]Album, error) {

	querySafeTerm := strings.ReplaceAll(term, " ", "+")
	url := fmt.Sprintf("%s/v1/catalog/%s/search?term=%s&types=albums", baseUrl, storefront, querySafeTerm)

	httpRes, err := a.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
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

func (a *AppleMusicClient) SearchSong(term string, storefront model.Market) ([]Song, error) {

	querySafeTerm := strings.ReplaceAll(term, " ", "+")
	url := fmt.Sprintf("%s/v1/catalog/%s/search?term=%s&types=songs", baseUrl, storefront, querySafeTerm)

	httpRes, err := a.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
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

func (a *AppleMusicClient) GetArtist(id string, storefront model.Market) (*Artist, error) {

	url := fmt.Sprintf("%s/v1/catalog/%s/artists/%s", baseUrl, storefront, id)

	httpRes, err := a.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
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

func (a *AppleMusicClient) GetAlbum(id string, storefront model.Market) (*Album, error) {

	url := fmt.Sprintf("%s/v1/catalog/%s/albums/%s?include=artists", baseUrl, storefront, id)

	httpRes, err := a.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
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

func (a *AppleMusicClient) GetSong(id string, storefront model.Market) (*Song, error) {

	url := fmt.Sprintf("%s/v1/catalog/%s/songs/%s?include=artists,albums", baseUrl, storefront, id)

	httpRes, err := a.client.Get(url)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, err
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
