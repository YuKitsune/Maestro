package appleMusic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// https://api.music.apple.com/v1/catalog/us/search

type appleMusicStreamingService struct {
	c                *http.Client
	shareLinkPattern *regexp.Regexp
}

func NewAppleMusicStreamingService(token string, shareLinkPattern string) streamingService.StreamingService {
	shareLinkPatternRegex := regexp.MustCompile(shareLinkPattern)
	return &appleMusicStreamingService{
		streamingService.NewClientWithBearerAuth(token),
		shareLinkPatternRegex,
	}
}

func (s *appleMusicStreamingService) Name() model.StreamingServiceKey {
	return "Apple Music"
}

func (s *appleMusicStreamingService) LinkBelongsToService(link string) bool {
	return s.shareLinkPattern.MatchString(link)
}

func (s *appleMusicStreamingService) SearchArtist(artist *model.Artist) (*model.Artist, error) {

	storefront := artist.GetMarket()
	term := strings.ReplaceAll(artist.Name, " ", "+")
	apiUrl := fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/search?term=%s&types=artists", storefront, term)

	httpRes, err := s.c.Get(apiUrl)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var apiRes *SearchResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return nil, err
	}

	if len(apiRes.Results.Artists.Data) == 0 {
		return nil, nil
	}

	appleMusicArtist := apiRes.Results.Artists.Data[0]
	link, err := url.Parse(appleMusicArtist.Attributes.Url)
	if err != nil {
		return nil, err
	}

	res := model.NewArtist(
		appleMusicArtist.Attributes.Name,
		nil, // Todo: artwork link
		s.Name(),
		storefront,
		link)
	return res, nil
}

func (s *appleMusicStreamingService) SearchAlbum(album *model.Album) (*model.Album, error) {

	storefront := album.GetMarket()
	term := strings.ReplaceAll(album.Name, " ", "+")
	apiUrl := fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/search?term=%s&types=albums", storefront, term)

	httpRes, err := s.c.Get(apiUrl)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var apiRes *SearchResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return nil, err
	}

	if len(apiRes.Results.Albums.Data) == 0 {
		return nil, nil
	}

	appleMusicAlbum := apiRes.Results.Albums.Data[0]
	artworkUrl, err := url.Parse(appleMusicAlbum.Attributes.Artwork.Url)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(appleMusicAlbum.Attributes.Url)
	if err != nil {
		return nil, err
	}

	res := model.NewAlbum(
		appleMusicAlbum.Attributes.Name,
		appleMusicAlbum.Attributes.ArtistName,
		artworkUrl,
		s.Name(),
		storefront,
		url)

	return res, nil
}

func (s *appleMusicStreamingService) SearchSong(song *model.Track) (*model.Track, error) {

	storefront := song.GetMarket()
	term := strings.ReplaceAll(song.Name, " ", "+")
	apiUrl := fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/search?term=%s&types=songs", storefront, term)

	httpRes, err := s.c.Get(apiUrl)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var apiRes *SearchResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return nil, err
	}

	if len(apiRes.Results.Songs.Data) == 0 {
		return nil, nil
	}

	appleMusicSong := apiRes.Results.Songs.Data[0]
	url, err := url.Parse(appleMusicSong.Attributes.Url)
	if err != nil {
		return nil, err
	}

	res := model.NewTrack(
		appleMusicSong.Attributes.Name,
		appleMusicSong.Attributes.ArtistName,
		appleMusicSong.Attributes.AlbumName,
		s.Name(),
		storefront,
		url)

	return res, nil
}

func (s *appleMusicStreamingService) SearchFromLink(link string) (model.Thing, error) {

	// example: https://music.apple.com/au/album/surrender/1585865534
	// format: 	https://music.apple.com/<storefront>/<artist|album|song>/<name>/<id>
	// name is irrelevant here, we only need the storefront and id

	matches := streamingService.FindStringSubmatchMap(s.shareLinkPattern, link)

	storefront := model.Market(matches["storefront"])
	typ := matches["type"]
	id := matches["id"]

	var apiUrl string
	var unmarshalFunc func(rb []byte) (model.Thing, error)

	switch typ {
	case "artist":
		apiUrl = fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/artists/%s", storefront, id)
		unmarshalFunc = func(rb []byte) (model.Thing, error) {
			var apiRes *ArtistsResult
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			if len(apiRes.Data) == 0 {
				return nil, fmt.Errorf("no %s found in region %s with id %s", typ, storefront, id)
			}

			foundArtist := apiRes.Data[0]

			link, err := url.Parse(foundArtist.Attributes.Url)
			if err != nil {
				return nil, err
			}

			artist := model.NewArtist(
				foundArtist.Attributes.Name,
				nil, // Todo: artwork link
				s.Name(),
				storefront,
				link)

			return artist, nil
		}

		break

	case "album":
		apiUrl = fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/albums/%s", storefront, id)
		unmarshalFunc = func(rb []byte) (model.Thing, error) {
			var apiRes *AlbumResult
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			foundAlbum := apiRes.Data[0]
			artworkUrl, err := url.Parse(foundAlbum.Attributes.Artwork.Url)
			if err != nil {
				return nil, err
			}

			url, err := url.Parse(foundAlbum.Attributes.Url)
			if err != nil {
				return nil, err
			}

			album := model.NewAlbum(
				foundAlbum.Attributes.Name,
				foundAlbum.Attributes.ArtistName,
				artworkUrl,
				s.Name(),
				storefront,
				url)

			return album, nil
		}
		break

	case "song":
		apiUrl = fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/songs/%s", storefront, id)
		unmarshalFunc = func(rb []byte) (model.Thing, error) {
			var apiRes *SongResult
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			foundSong := apiRes.Data[0]
			url, err := url.Parse(foundSong.Attributes.Url)
			if err != nil {
				return nil, err
			}

			song := model.NewTrack(
				foundSong.Attributes.Name,
				foundSong.Attributes.ArtistName,
				foundSong.Attributes.AlbumName,
				s.Name(),
				storefront,
				url)

			return song, nil
		}
		break

	default:
		return nil, fmt.Errorf("unknown type %s", typ)
	}

	httpRes, err := s.c.Get(apiUrl)
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no %s found in region %s with id %s", typ, storefront, id)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	res, err := unmarshalFunc(resBytes)

	return res, err
}

