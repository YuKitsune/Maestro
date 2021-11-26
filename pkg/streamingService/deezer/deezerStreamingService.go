package deezer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"maestro/pkg/streamingService"
	"net/http"
	"net/url"
	"regexp"
)

type deezerStreamingService struct {
	client *http.Client
	shareLinkPattern *regexp.Regexp
}

func getActualLink(link string, linkRegexp *regexp.Regexp) (string, error) {
	var actualLink string

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			urlString := req.URL.String()
			if linkRegexp.MatchString(urlString) {
				actualLink = urlString
				return http.ErrUseLastResponse
			}

			return nil
		},
	}

	_, err := client.Get(link)
	return actualLink, err
}

func NewDeezerStreamingService(shareLinkPattern string) streamingService.StreamingService {
	shareLinkPatternRegex := regexp.MustCompile(shareLinkPattern)
	return &deezerStreamingService{&http.Client{}, shareLinkPatternRegex}
}

func (s *deezerStreamingService) LinkBelongsToService(link string) bool {
	return s.shareLinkPattern.MatchString(link)
}

func (s *deezerStreamingService) Name() string {
	return "Deezer"
}

func (s *deezerStreamingService) SearchArtist(artist *streamingService.Artist) (*streamingService.Artist, error) {

	q := url.QueryEscape(fmt.Sprintf("artist:\"%s\"", artist.Name))
	apiUrl := fmt.Sprintf("https://api.deezer.com/search/artist?q=%s", q)

	httpRes, err := s.client.Get(apiUrl)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var apiRes *searchArtistResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return nil, err
	}

	if len(apiRes.Data) == 0 {
		return nil, nil
	}

	deezerArtist := apiRes.Data[0]
	return &streamingService.Artist{
		Name:       deezerArtist.Name,
		ArtworkUrl: deezerArtist.Picture,
		Url:        deezerArtist.Link,
	}, nil
}

func (s *deezerStreamingService) SearchAlbum(album *streamingService.Album) (*streamingService.Album, error) {

	q := url.QueryEscape(fmt.Sprintf("artist:\"%s\" album:\"%s\"", album.ArtistName, album.Name))
	apiUrl := fmt.Sprintf("https://api.deezer.com/search/album?q=%s", q)

	httpRes, err := s.client.Get(apiUrl)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var apiRes *searchAlbumResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return nil, err
	}

	if len(apiRes.Data) == 0 {
		return nil, nil
	}

	deezerAlbum :=  apiRes.Data[0]
	return &streamingService.Album{
		Name:       deezerAlbum.Title,
		ArtistName: deezerAlbum.Artist.Name,
		ArtworkUrl: deezerAlbum.Cover,
		Url:        deezerAlbum.Link,
	}, nil
}

func (s *deezerStreamingService) SearchSong(song *streamingService.Song) (*streamingService.Song, error) {

	q := url.QueryEscape(fmt.Sprintf("artist:\"%s\" album:\"%s\" track:\"%s\"", song.ArtistName, song.AlbumName, song.Name))
	apiUrl := fmt.Sprintf("https://api.deezer.com/search/track?q=%s", q)

	httpRes, err := s.client.Get(apiUrl)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var apiRes *searchTrackResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return nil, err
	}

	if len(apiRes.Data) == 0 {
		return nil, nil
	}

	deezerTrack :=  apiRes.Data[0]
	return &streamingService.Song{
		Name:       deezerTrack.Title,
		ArtistName: deezerTrack.Artist.Name,
		AlbumName:  deezerTrack.Album.Title,
		Url:        deezerTrack.Link,
	}, nil
}

func (s *deezerStreamingService) SearchFromLink(link string) (streamingService.Thing, error) {

	// Share link: https://deezer.page.link/szbWkX6rKbfJ8XCD6
	// This goes through some redirects until we get to here:
	// example: https://www.deezer.com/en/track/606334862<some stuff i don't care about>
	// format: 	https://www.deezer.com/<lang>/<artist|album|track>/<id>
	// Todo: How we gonna get the region?

	actualLink, err := getActualLink(link, s.shareLinkPattern)
	if err != nil {
		return nil, err
	}

	matches := findStringSubmatchMap(s.shareLinkPattern, actualLink)

	// store := matches["storefront"]
	typ := matches["type"]
	id := matches["id"]

	var url string
	var unmarshalFunc func(rb []byte) (streamingService.Thing, error)

	switch typ {
	case "artist":
		url = fmt.Sprintf("https://api.deezer.com/artist/%s", id)
		unmarshalFunc = func (rb []byte) (streamingService.Thing, error) {
			var apiRes *Artist
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			artist := &streamingService.Artist{
				Name:       apiRes.Name,
				Url:        apiRes.Link,
				ArtworkUrl: apiRes.Picture,
			}

			return artist, nil
		}

		break

	case "album":
		url = fmt.Sprintf("https://api.deezer.com/album/%s", id)
		unmarshalFunc = func (rb []byte) (streamingService.Thing, error) {
			var apiRes *Album
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			album := &streamingService.Album{
				Name:       apiRes.Title,
				ArtistName: apiRes.Artist.Name,
				ArtworkUrl: apiRes.Cover,
				Url:        apiRes.Link,
			}

			return album, nil
		}
		break

	case "track":
		url = fmt.Sprintf("https://api.deezer.com/track/%s", id)
		unmarshalFunc = func (rb []byte) (streamingService.Thing, error) {
			var apiRes *Track
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			song := &streamingService.Song{
				Name:        apiRes.Title,
				ArtistName:  apiRes.Artist.Name,
				AlbumName:   apiRes.Album.Title,
				Url:         apiRes.Link,
			}

			return song, nil
		}
		break

	default:
		return nil, fmt.Errorf("unknown type %s", typ)
	}

	httpRes, err := s.client.Get(url)
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusNotFound {
		// return nil, fmt.Errorf("no %s found in region %s with id %s", typ, store, id)
		return nil, fmt.Errorf("no %s found with id %s", typ, id)
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