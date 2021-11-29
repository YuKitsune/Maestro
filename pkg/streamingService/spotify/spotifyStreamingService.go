package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/zmb3/spotify"
	"io/ioutil"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"net/url"
	"regexp"
)

type spotifyStreamingService struct {
	client           *spotify.Client
	shareLinkPattern *regexp.Regexp
}

func GetAccessToken(clientId string, secret string) (token string, error error) {
	tokenUrl := "https://accounts.spotify.com/api/token"

	reqToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientId, secret)))
	client := streamingService.NewClientWithBasicAuth(reqToken)

	res, err := client.PostForm(tokenUrl, url.Values{
		"grant_type": {"client_credentials"},
	})
	defer res.Body.Close()

	if err != nil {
		return token, err
	}

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return token, err
	}

	var resMap map[string]interface{}
	err = json.Unmarshal(resBytes, &resMap)
	if err != nil {
		return token, err
	}

	token = resMap["access_token"].(string)

	return token, nil
}

func NewSpotifyStreamingService(clientId string, clientSecret string, shareLinkPattern string) (streamingService.StreamingService, error) {
	shareLinkPatternRegex := regexp.MustCompile(shareLinkPattern)

	token, err := GetAccessToken(clientId, clientSecret)
	if err != nil {
		return nil, err
	}

	c := streamingService.NewClientWithBearerAuth(token)
	sc := spotify.NewClient(c)
	return &spotifyStreamingService{&sc, shareLinkPatternRegex}, nil
}

func (s *spotifyStreamingService) Name() model.StreamingServiceKey {
	return "Spotify"
}

func (s *spotifyStreamingService) LinkBelongsToService(link string) bool {
	return s.shareLinkPattern.MatchString(link)
}

func (s *spotifyStreamingService) SearchArtist(artist *model.Artist) (*model.Artist, error) {

	country := artist.Market.String()

	q := fmt.Sprintf("artist:\"%s\"", artist.Name)
	searchRes, err := s.client.SearchOpt(q, spotify.SearchTypeArtist, &spotify.Options{
		Country: &country,
	})
	if err != nil {
		return nil, err
	}

	if searchRes.Artists == nil || len(searchRes.Artists.Artists) == 0 {
		return nil, nil
	}

	spotifyArtist := searchRes.Artists.Artists[0]
	res := model.NewArtist(
		spotifyArtist.Name,
		imageUrl(spotifyArtist.Images),
		s.Name(),
		model.DefaultMarket,
		spotifyArtist.ExternalURLs["spotify"])

	return res, nil
}

func (s *spotifyStreamingService) SearchAlbum(album *model.Album) (*model.Album, error) {

	country := album.Market.String()

	q := fmt.Sprintf("artist:\"%s\" album:\"%s\"", album.ArtistNames[0], album.Name)
	searchRes, err := s.client.SearchOpt(q, spotify.SearchTypeAlbum, &spotify.Options{
		Country: &country,
	})
	if err != nil {
		return nil, err
	}

	if searchRes.Albums == nil || len(searchRes.Albums.Albums) == 0 {
		return nil, nil
	}

	spotifyAlbum := searchRes.Albums.Albums[0]
	res := model.NewAlbum(
		spotifyAlbum.Name,
		artistName(spotifyAlbum.Artists),
		imageUrl(spotifyAlbum.Images),
		s.Name(),
		model.DefaultMarket,
		spotifyAlbum.ExternalURLs["spotify"])

	return res, nil
}

func (s *spotifyStreamingService) SearchSong(track *model.Track) (*model.Track, error) {

	country := track.Market.String()

	q := fmt.Sprintf("artist:\"%s\" album:\"%s\" track:\"%s\"", track.ArtistNames[0], track.AlbumName, track.Name)
	searchRes, err := s.client.SearchOpt(q, spotify.SearchTypeTrack, &spotify.Options{
		Country: &country,
	})
	if err != nil {
		return nil, err
	}

	if searchRes.Tracks == nil || len(searchRes.Tracks.Tracks) == 0 {
		return nil, nil
	}

	spotifyTrack := searchRes.Tracks.Tracks[0]
	res := model.NewTrack(
		spotifyTrack.Name,
		artistName(spotifyTrack.Artists),
		spotifyTrack.Album.Name,
		s.Name(),
		model.DefaultMarket,
		spotifyTrack.ExternalURLs["spotify"])

	return res, nil
}

func (s *spotifyStreamingService) SearchFromLink(link string) (model.Thing, error) {

	// example: https://open.spotify.com/track/4cOdK2wGLETKBW3PvgPWqT?si=10587ef152a8493f
	// format: 	https://open.spotify.com/<artist|album|track>/<id>?si=<user specific token that i don't care about>
	// Todo: How we gonna get the region?

	matches := findStringSubmatchMap(s.shareLinkPattern, link)

	// region := matches["region"]
	typ := matches["type"]
	id := spotify.ID(matches["id"])

	switch typ {
	case "artist":
		foundArtist, err := s.client.GetArtist(id)
		if err != nil {
			return nil, err
		}

		artist := model.NewArtist(
			foundArtist.Name,
			imageUrl(foundArtist.Images),
			s.Name(),
			model.DefaultMarket,
			foundArtist.ExternalURLs["spotify"])

		return artist, nil

	case "album":
		foundAlbum, err := s.client.GetAlbum(id)
		if err != nil {
			return nil, err
		}

		album := model.NewAlbum(
			foundAlbum.Name,
			artistName(foundAlbum.Artists),
			imageUrl(foundAlbum.Images),
			s.Name(),
			model.DefaultMarket,
			foundAlbum.ExternalURLs["spotify"])

		return album, nil

	case "track":
		foundTrack, err := s.client.GetTrack(id)
		if err != nil {
			return nil, err
		}

		track := model.NewTrack(
			foundTrack.Name,
			artistName(foundTrack.Artists),
			foundTrack.Album.Name,
			s.Name(),
			model.DefaultMarket,
			foundTrack.ExternalURLs["spotify"])

		return track, nil

	default:
		return nil, fmt.Errorf("unknown type %s", typ)
	}
}

func artistName(artists []spotify.SimpleArtist) []string {

	var names []string
	if len(artists) > 0 {
		for _, artist := range artists {
			names = append(names, artist.Name)
		}
	}

	return names
}

func imageUrl(imgs []spotify.Image) string {
	if len(imgs) > 0 {
		return imgs[0].URL
	}

	return ""
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
