package appleMusic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"maestro/pkg/streamingService"
	"net/http"
	"regexp"
	"strings"
)

// https://api.music.apple.com/v1/catalog/us/search

type appleMusicStreamingService struct {
	c *http.Client
	shareLinkPattern *regexp.Regexp
}

func NewAppleMusicStreamingService(token string, shareLinkPattern string) streamingService.StreamingService {
	shareLinkPatternRegex := regexp.MustCompile(shareLinkPattern)
	return &appleMusicStreamingService{
		streamingService.NewClientWithBearerAuth(token),
		shareLinkPatternRegex,
	}
}

func (s *appleMusicStreamingService) Name() string {
	return "Apple Music"
}

func (s *appleMusicStreamingService) LinkBelongsToService(link string) bool {
	return s.shareLinkPattern.MatchString(link)
}

func (s *appleMusicStreamingService) SearchArtist(artist *streamingService.Artist) (*streamingService.Artist, error) {

	region := artist.GetRegion()
	term := strings.ReplaceAll(artist.Name, " ", "+")
	url := fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/search?term=%s&types=artists", region, term)

	httpRes, err := s.c.Get(url)
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
	return &streamingService.Artist{
		Name:   appleMusicArtist.Attributes.Name,
		Url:    appleMusicArtist.Attributes.Url,
	}, nil
}

func (s *appleMusicStreamingService) SearchAlbum(album *streamingService.Album) (*streamingService.Album, error) {

	region := album.GetRegion()
	term := strings.ReplaceAll(album.Name, " ", "+")
	url := fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/search?term=%s&types=albums", region, term)

	httpRes, err := s.c.Get(url)
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
	return &streamingService.Album{
		Name:       appleMusicAlbum.Attributes.Name,
		ArtistName: appleMusicAlbum.Attributes.ArtistName,
		ArtworkUrl: appleMusicAlbum.Attributes.Artwork.Url,
		Url:        appleMusicAlbum.Attributes.Url,
	}, nil
}

func (s *appleMusicStreamingService) SearchSong(song *streamingService.Song) (*streamingService.Song, error) {

	region := song.GetRegion()
	term := strings.ReplaceAll(song.Name, " ", "+")
	url := fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/search?term=%s&types=songs", region, term)

	httpRes, err := s.c.Get(url)
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
	return &streamingService.Song{
		Name:       appleMusicSong.Attributes.Name,
		ArtistName: appleMusicSong.Attributes.ArtistName,
		AlbumName:  appleMusicSong.Attributes.AlbumName,
		Url:        appleMusicSong.Attributes.Url,
	}, nil
}

func (s *appleMusicStreamingService) SearchFromLink(link string) (streamingService.Thing, error) {

	// example: https://music.apple.com/au/album/surrender/1585865534
	// format: 	https://music.apple.com/<storefront>/<artist|album|song>/<name>/<id>
	// name is irrelevant here, we only need the storefront and id

	matches := findStringSubmatchMap(s.shareLinkPattern, link)

	store := matches["storefront"]
	typ := matches["type"]
	id := matches["id"]

	var url string
	var unmarshalFunc func(rb []byte) (streamingService.Thing, error)

	switch typ {
	case "artist":
		url = fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/artists/%s", store, id)
		unmarshalFunc = func (rb []byte) (streamingService.Thing, error) {
			var apiRes *ArtistsResult
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			if len(apiRes.Data) == 0 {
				return nil, fmt.Errorf("no %s found in region %s with id %s", typ, store, id)
			}

			foundArtist := apiRes.Data[0]
			artist := &streamingService.Artist{
				Name:       foundArtist.Attributes.Name,
				Url:        foundArtist.Attributes.Url,
				ArtworkUrl: "",
			}

			return artist, nil
		}

		break

	case "album":
		url = fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/albums/%s", store, id)
		unmarshalFunc = func (rb []byte) (streamingService.Thing, error) {
			var apiRes *AlbumResult
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			foundAlbum := apiRes.Data[0]
			album := &streamingService.Album{
				Name:       foundAlbum.Attributes.Name,
				ArtistName: foundAlbum.Attributes.ArtistName,
				ArtworkUrl: foundAlbum.Attributes.Artwork.Url,
				Url:        foundAlbum.Attributes.Url,
			}

			return album, nil
		}
		break

	case "song":
		url = fmt.Sprintf("https://api.music.apple.com/v1/catalog/%s/songs/%s", store, id)
		unmarshalFunc = func (rb []byte) (streamingService.Thing, error) {
			var apiRes *SongResult
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			foundSong := apiRes.Data[0]
			song := &streamingService.Song{
				Name:        foundSong.Attributes.Name,
				ArtistName:  foundSong.Attributes.ArtistName,
				AlbumName:   foundSong.Attributes.AlbumName,
				Url:         foundSong.Attributes.Url,
			}

			return song, nil
		}
		break

	default:
		return nil, fmt.Errorf("unknown type %s", typ)
	}

	httpRes, err := s.c.Get(url)
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no %s found in region %s with id %s", typ, store, id)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	res, err := unmarshalFunc(resBytes)

	return res, err
}

func findStringSubmatchMap(r *regexp.Regexp, s string) map[string]string {

	matches := r.FindStringSubmatch(s)
	names := r.SubexpNames()

	result := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}

	return result
}
