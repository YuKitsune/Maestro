package deezer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"net/http"
	"net/url"
	"regexp"
)

type deezerStreamingService struct {
	client           *http.Client
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

func (s *deezerStreamingService) Name() model.StreamingServiceKey {
	return "Deezer"
}

func (s *deezerStreamingService) SearchArtist(artist *model.Artist) (*model.Artist, error) {

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
	url, err := url.Parse(deezerArtist.Link)
	if err != nil {
		return nil, err
	}

	artUrl, err := url.Parse(deezerArtist.Picture)
	if err != nil {
		return nil, err
	}

	res := model.NewArtist(
		deezerArtist.Name,
		artUrl,
		s.Name(),
		model.DefaultMarket,
		url)

	return res, nil
}

func (s *deezerStreamingService) SearchAlbum(album *model.Album) (*model.Album, error) {

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

	deezerAlbum := apiRes.Data[0]
	url, err := url.Parse(deezerAlbum.Link)
	if err != nil {
		return nil, err
	}

	artUrl, err := url.Parse(deezerAlbum.Cover)
	if err != nil {
		return nil, err
	}

	res := model.NewAlbum(
		deezerAlbum.Title,
		deezerAlbum.Artist.Name,
		artUrl,
		s.Name(),
		model.DefaultMarket,
		url)

	return res, nil
}

func (s *deezerStreamingService) SearchSong(song *model.Track) (*model.Track, error) {

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

	deezerTrack := apiRes.Data[0]
	url, err := url.Parse(deezerTrack.Link)
	if err != nil {
		return nil, err
	}

	res := model.NewTrack(
		deezerTrack.Title,
		deezerTrack.Artist.Name,
		deezerTrack.Album.Title,
		s.Name(),
		model.DefaultMarket,
		url)

	return res, nil
}

func (s *deezerStreamingService) SearchFromLink(link string) (model.Thing, error) {

	// Share link: https://deezer.page.link/szbWkX6rKbfJ8XCD6
	// This goes through some redirects until we get to here:
	// example: https://www.deezer.com/en/track/606334862<some stuff i don't care about>
	// format: 	https://www.deezer.com/<lang>/<artist|album|track>/<id>
	// Todo: How we gonna get the region?

	actualLink, err := getActualLink(link, s.shareLinkPattern)
	if err != nil {
		return nil, err
	}

	matches := streamingService.FindStringSubmatchMap(s.shareLinkPattern, actualLink)

	// store := matches["storefront"]
	typ := matches["type"]
	id := matches["id"]

	var apiUrl string
	var unmarshalFunc func(rb []byte) (model.Thing, error)

	switch typ {
	case "artist":
		apiUrl = fmt.Sprintf("https://api.deezer.com/artist/%s", id)
		unmarshalFunc = func(rb []byte) (model.Thing, error) {
			var apiRes *Artist
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			url, err := url.Parse(apiRes.Link)
			if err != nil {
				return nil, err
			}

			artUrl, err := url.Parse(apiRes.Picture)
			if err != nil {
				return nil, err
			}

			artist := model.NewArtist(
				apiRes.Name,
				artUrl,
				s.Name(),
				model.DefaultMarket,
				url)
			return artist, nil
		}

		break

	case "album":
		apiUrl = fmt.Sprintf("https://api.deezer.com/album/%s", id)
		unmarshalFunc = func(rb []byte) (model.Thing, error) {
			var apiRes *Album
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			url, err := url.Parse(apiRes.Link)
			if err != nil {
				return nil, err
			}

			artUrl, err := url.Parse(apiRes.Cover)
			if err != nil {
				return nil, err
			}

			album := model.NewAlbum(
				apiRes.Title,
				apiRes.Artist.Name,
				artUrl,
				s.Name(),
				model.DefaultMarket,
				url)

			return album, nil
		}
		break

	case "track":
		apiUrl = fmt.Sprintf("https://api.deezer.com/track/%s", id)
		unmarshalFunc = func(rb []byte) (model.Thing, error) {
			var apiRes *Track
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			url, err := url.Parse(apiRes.Link)
			if err != nil {
				return nil, err
			}

			track := model.NewTrack(
				apiRes.Title,
				apiRes.Artist.Name,
				apiRes.Album.Title,
				s.Name(),
				model.DefaultMarket,
				url)

			return track, nil
		}
		break

	default:
		return nil, fmt.Errorf("unknown type %s", typ)
	}

	httpRes, err := s.client.Get(apiUrl)
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
