package deezer

import (
	"fmt"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"net/http"
	"regexp"
)

type deezerStreamingService struct {
	client           *DeezerClient
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
	return &deezerStreamingService{NewDeezerClient(), shareLinkPatternRegex}
}

func (s *deezerStreamingService) LinkBelongsToService(link string) bool {
	return s.shareLinkPattern.MatchString(link)
}

func (s *deezerStreamingService) Name() model.StreamingServiceKey {
	return "Deezer"
}

func (s *deezerStreamingService) SearchArtist(artist *model.Artist) (*model.Artist, error) {

	searchRes, err := s.client.SearchArtist(artist.Name)
	if err != nil {
		return nil, err
	}

	// Todo: Narrow down results
	deezerArtist := searchRes[0]

	res := model.NewArtist(
		deezerArtist.Name,
		deezerArtist.Picture,
		s.Name(),
		model.DefaultMarket,
		deezerArtist.Link)

	return res, nil
}

func (s *deezerStreamingService) SearchAlbum(album *model.Album) (*model.Album, error) {

	// Deezer only has one artist per track/album, need to check each artist

	var res *model.Album
	for _, artistName := range album.ArtistNames {

		searchRes, err := s.client.SearchAlbum(artistName, album.Name)
		if err != nil {
			return nil, err
		}

		if searchRes == nil || len(searchRes) == 0 {
			continue
		}

		// Todo: Narrow down results
		deezerAlbum := searchRes[0]
		res = model.NewAlbum(
			deezerAlbum.Title,
			[]string {deezerAlbum.Artist.Name},
			deezerAlbum.Cover,
			s.Name(),
			model.DefaultMarket,
			deezerAlbum.Link)
	}

	return res, nil
}

func (s *deezerStreamingService) SearchSong(track *model.Track) (*model.Track, error) {

	var res *model.Track
	for _, artistName := range track.ArtistNames {

		searchRes, err := s.client.SearchTrack(artistName, track.AlbumName, track.Name)
		if err != nil {
			return nil, err
		}

		if searchRes == nil || len(searchRes) == 0 {
			continue
		}

		deezerTrack := searchRes[0]
		res = model.NewTrack(
			deezerTrack.Title,
			[]string {deezerTrack.Artist.Name},
			deezerTrack.Album.Title,
			s.Name(),
			model.DefaultMarket,
			deezerTrack.Link)
	}

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

	switch typ {
	case "artist":
		foundArtist, err := s.client.GetArtist(id)
		if err != nil {
			return nil, err
		}

		artist := model.NewArtist(
			foundArtist.Name,
			foundArtist.Picture,
			s.Name(),
			model.DefaultMarket,
			foundArtist.Link)

		return artist, nil

	case "album":
		foundAlbum, err := s.client.GetAlbum(id)
		if err != nil {
			return nil, err
		}

		album := model.NewAlbum(
			foundAlbum.Title,
			[]string {foundAlbum.Artist.Name}, // Todo:
			foundAlbum.Cover,
			s.Name(),
			model.DefaultMarket,
			foundAlbum.Link)

		return album, nil

	case "track":
		foundTrack, err := s.client.GetTrack(id)
		if err != nil {
			return nil, err
		}

		track := model.NewTrack(
			foundTrack.Title,
			[]string {foundTrack.Artist.Name}, // Todo:
			foundTrack.Album.Title,
			s.Name(),
			model.DefaultMarket,
			foundTrack.Link)

		return track, nil

	default:
		return nil, fmt.Errorf("unknown type %s", typ)
	}
}

func (s *deezerStreamingService) CleanLink(link string) string {

	match := s.shareLinkPattern.FindStringIndex(link)
	if len(match) > 0 {
		return link[match[0]:match[1]]
	}

	return link
}
