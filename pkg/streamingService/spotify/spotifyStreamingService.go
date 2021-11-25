package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/zmb3/spotify"
	"io/ioutil"
	"maestro/pkg/streamingService"
	"net/url"
	"regexp"
)

type spotifyStreamingService struct {
	client *spotify.Client
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

func NewSpotifyStreamingService(clientId string, clientSecret string) (streamingService.StreamingService, error) {

	token, err := GetAccessToken(clientId, clientSecret)
	if err != nil {
		return nil, err
	}

	c := streamingService.NewClientWithBearerAuth(token)
	sc := spotify.NewClient(c)
	return &spotifyStreamingService{client: &sc}, nil
}

func (s *spotifyStreamingService) Name() string {
	return "Spotify"
}

func (s *spotifyStreamingService) SearchArtist(name string, region streamingService.Region) (res []streamingService.Artist, err error) {

	country := streamingService.RegionToString(region)
	so := &spotify.Options{
		Country:   &country,
	}

	searchRes, err := s.client.SearchOpt(name, spotify.SearchTypeArtist, so)
	if err != nil {
		return res, err
	}

	for _, spotifyArtist := range searchRes.Artists.Artists {

		var imageUrl string
		if len(spotifyArtist.Images) > 0 {
			imageUrl = spotifyArtist.Images[0].URL
		}

		url := spotifyArtist.ExternalURLs["spotify"]
		artist := streamingService.Artist{
			Name:       spotifyArtist.Name,
			ArtworkUrl: imageUrl,
			Url:        url,
		}

		res = append(res, artist)
	}

	return res, nil
}

func (s *spotifyStreamingService) SearchAlbum(name string, region streamingService.Region) (res []streamingService.Album, err error) {

	country := streamingService.RegionToString(region)
	so := &spotify.Options{
		Country:   &country,
	}

	searchRes, err := s.client.SearchOpt(name, spotify.SearchTypeArtist, so)
	if err != nil {
		return res, err
	}

	for _, spotifyAlbum := range searchRes.Albums.Albums {

		url := spotifyAlbum.ExternalURLs["spotify"]

		var imageUrl string
		if len(spotifyAlbum.Images) > 0 {
			imageUrl = spotifyAlbum.Images[0].URL
		}

		album := streamingService.Album{
			Name:       spotifyAlbum.Name,
			ArtistName: artistName(spotifyAlbum.Artists),
			ArtworkUrl: imageUrl,
			Url:        url,
		}

		res = append(res, album)
	}

	return res, nil
}

func (s *spotifyStreamingService) SearchSong(name string, region streamingService.Region) (res []streamingService.Song, err error) {

	country := streamingService.RegionToString(region)
	so := &spotify.Options{
		Country:   &country,
	}

	searchRes, err := s.client.SearchOpt(name, spotify.SearchTypeArtist, so)
	if err != nil {
		return res, err
	}

	for _, spotifySong := range searchRes.Tracks.Tracks {

		url := spotifySong.ExternalURLs["spotify"]

		song := streamingService.Song{
			Name:       spotifySong.Name,
			ArtistName: artistName(spotifySong.Artists),
			AlbumName:  spotifySong.Album.Name,
			Url:        url,
		}

		res = append(res, song)
	}

	return res, nil
}

func (s *spotifyStreamingService) SearchFromLink(link string) (streamingService.Thing, error) {

	// example: https://open.spotify.com/track/4cOdK2wGLETKBW3PvgPWqT?si=10587ef152a8493f
	// format: 	https://open.spotify.com/<artist|album|track>/<id>?si=<user specific token that i don't care about>
	// Todo: How we gonna get the region?

	// Todo: Move pattern to config
	pattern := "(?:https:\\/\\/open\\.spotify\\.com\\/)(?P<type>[A-Za-z]+)\\/(?P<id>[A-Za-z0-9]+).*"
	linkRegexp := regexp.MustCompile(pattern)

	matches := findStringSubmatchMap(linkRegexp, link)

	// region := matches["region"]
	typ := matches["type"]
	id := spotify.ID(matches["id"])

	switch typ {
	case "artist":
		foundArtist, err := s.client.GetArtist(id)
		if err != nil {
			return nil, err
		}

		artist := &streamingService.Artist{
			Name:       foundArtist.Name,
			ArtworkUrl: imageUrl(foundArtist.Images),
			Url:        foundArtist.ExternalURLs["spotify"],
		}

		return artist, nil

	case "album":
		foundAlbum, err := s.client.GetAlbum(id)
		if err != nil {
			return nil, err
		}

		album := &streamingService.Artist{
			Name:       foundAlbum.Name,
			ArtworkUrl: imageUrl(foundAlbum.Images),
			Url:        foundAlbum.ExternalURLs["spotify"],
		}

		return album, nil

	case "track":
		foundTrack, err := s.client.GetTrack(id)
		if err != nil {
			return nil, err
		}

		track := &streamingService.Song{
			Name:       foundTrack.Name,
			ArtistName: artistName(foundTrack.Artists),
			AlbumName:  foundTrack.Album.Name,
			Url:        foundTrack.ExternalURLs["spotify"],
		}

		return track, nil

	default:
		return nil, fmt.Errorf("unknown type %s", typ)
	}
}

func artistName(artists []spotify.SimpleArtist) string {

	var name string
	if len(artists) > 0 {
		for i, artist := range artists {
			if i > 0 && i == len(artists)-1 {
				name += ", "
			}

			name += artist.Name
		}
	}

	return name
}

func imageUrl(imgs []spotify.Image) string {
	var url string
	if len(imgs) > 0 {
		url = imgs[0].URL
	}

	return url
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
