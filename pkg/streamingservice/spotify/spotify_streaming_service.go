package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/yukitsune/maestro/pkg/metrics"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingservice"
	"github.com/zmb3/spotify"
	"io/ioutil"
	"net/url"
	"regexp"
)

// Links:
//   Artist: https://open.spotify.com/artist/3dz0NnIZhtKKeXZxLOxCam?si=59MxgP9hT72oENCW_7Pi8g
//    Album: https://open.spotify.com/album/4Hjqdhj5rh816i1dfcUEaM?si=Dc3aTkvPT22eZA6viIkdKA
//    Track: https://open.spotify.com/track/36kCSJg8ZBwiSCUECFKGUy?si=df04501fab5540d2
// Playlist: https://open.spotify.com/playlist/6CmKR9e3DTdyNdfr4NBeSG?si=c716a24afbf7402d
const shareLinkRegex = "(https?:\\/\\/)?open\\.spotify\\.com\\/(?P<type>[A-Za-z]+)\\/(?P<id>[A-Za-z0-9]+)"

type spotifyStreamingService struct {
	client           *spotify.Client
	shareLinkPattern *regexp.Regexp
	metricsRecorder  metrics.Recorder
}

func GetAccessToken(clientID string, secret string) (token string, error error) {
	tokenURL := "https://accounts.spotify.com/api/token"

	reqToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, secret)))
	client := streamingservice.NewClientWithBasicAuth(reqToken)

	res, err := client.PostForm(tokenURL, url.Values{
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

	errorMessage, ok := resMap["error"].(string)
	if ok && len(errorMessage) > 0 {
		return "", fmt.Errorf("failed to get access token: %s", errorMessage)
	}

	token, ok = resMap["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("could not find access token in response")
	}

	return token, nil
}

func NewSpotifyStreamingService(cfg *Config, mr metrics.Recorder) (streamingservice.StreamingService, error) {
	shareLinkPatternRegex := regexp.MustCompile(shareLinkRegex)

	go mr.CountSpotifyRequest()
	token, err := GetAccessToken(cfg.ClientID, cfg.ClientSecret)
	if err != nil {
		return nil, err
	}

	c := streamingservice.NewClientWithBearerAuth(token)
	sc := spotify.NewClient(c)
	return &spotifyStreamingService{
		&sc,
		shareLinkPatternRegex,
		mr,
	}, nil
}

func (s *spotifyStreamingService) Key() model.StreamingServiceKey {
	return Key
}

func (s *spotifyStreamingService) LinkBelongsToService(link string) bool {
	return s.shareLinkPattern.MatchString(link)
}

func (s *spotifyStreamingService) SearchArtist(artist *model.Artist) (*model.Artist, bool, error) {

	country := artist.Market.String()

	go s.metricsRecorder.CountSpotifyRequest()

	q := fmt.Sprintf("artist:\"%s\"", artist.Name)
	searchRes, err := s.client.SearchOpt(q, spotify.SearchTypeArtist, &spotify.Options{
		Country: &country,
	})
	if err != nil {
		return nil, false, err
	}

	if searchRes.Artists == nil || len(searchRes.Artists.Artists) == 0 {
		return nil, false, nil
	}

	spotifyArtist := searchRes.Artists.Artists[0]
	res := model.NewArtist(
		spotifyArtist.Name,
		imageURL(spotifyArtist.Images),
		s.Key(),
		model.DefaultMarket,
		spotifyArtist.ExternalURLs["spotify"])

	return res, true, nil
}

func (s *spotifyStreamingService) SearchAlbum(album *model.Album) (*model.Album, bool, error) {

	country := album.Market.String()

	// Spotify search API doesn't like multiple artist names in the search query
	// need to query each artist separately
	// Sigh...

	var res *model.Album
	for _, name := range album.ArtistNames {

		go s.metricsRecorder.CountSpotifyRequest()

		q := fmt.Sprintf("artist:\"%s\" album:\"%s\"", name, album.Name)
		searchRes, err := s.client.SearchOpt(q, spotify.SearchTypeAlbum, &spotify.Options{
			Country: &country,
		})
		if err != nil {
			return nil, false, err
		}

		if searchRes.Albums == nil || len(searchRes.Albums.Albums) == 0 {
			continue
		}

		// Todo: Narrow down results
		spotifyAlbum := searchRes.Albums.Albums[0]

		res = model.NewAlbum(
			spotifyAlbum.Name,
			artistName(spotifyAlbum.Artists),
			imageURL(spotifyAlbum.Images),
			s.Key(),
			model.DefaultMarket,
			spotifyAlbum.ExternalURLs["spotify"])
	}

	return res, res != nil, nil
}

func (s *spotifyStreamingService) GetTrackByIsrc(isrc string) (*model.Track, bool, error) {

	q := fmt.Sprintf("isrc:\"%s\"", isrc)

	searchRes, err := s.client.SearchOpt(q, spotify.SearchTypeTrack, &spotify.Options{})
	if err != nil {
		return nil, false, err
	}

	if err != nil {
		return nil, false, err
	}

	if searchRes.Tracks == nil || len(searchRes.Tracks.Tracks) == 0 {
		return nil, false, nil
	}

	// Todo: Narrow down results
	spotifyTrack := searchRes.Tracks.Tracks[0]

	res := model.NewTrack(
		spotifyTrack.ExternalIDs["isrc"],
		spotifyTrack.Name,
		artistName(spotifyTrack.Artists),
		spotifyTrack.Album.Name,
		imageURL(spotifyTrack.Album.Images),
		s.Key(),
		model.DefaultMarket,
		spotifyTrack.ExternalURLs["spotify"])

	return res, true, nil
}

func (s *spotifyStreamingService) SearchTrack(track *model.Track) (*model.Track, bool, error) {

	country := track.Market.String()

	// Spotify search API doesn't like multiple artist names in the search query
	// need to query each artist separately
	// Sigh...

	var res *model.Track
	for _, name := range track.ArtistNames {

		var spotifyTrack spotify.FullTrack
		var q string
		if len(track.Isrc) > 0 {
			q = fmt.Sprintf("isrc:\"%s\"", track.Isrc)
		} else {
			q = fmt.Sprintf("artist:\"%s\" album:\"%s\" track:\"%s\"", name, track.AlbumName, track.Name)
		}

		go s.metricsRecorder.CountSpotifyRequest()

		searchRes, err := s.client.SearchOpt(q, spotify.SearchTypeTrack, &spotify.Options{
			Country: &country,
		})
		if err != nil {
			return nil, false, err
		}

		if searchRes.Tracks == nil || len(searchRes.Tracks.Tracks) == 0 {
			continue
		}

		// Todo: Narrow down results
		spotifyTrack = searchRes.Tracks.Tracks[0]

		res = model.NewTrack(
			spotifyTrack.ExternalIDs["isrc"],
			spotifyTrack.Name,
			artistName(spotifyTrack.Artists),
			spotifyTrack.Album.Name,
			imageURL(spotifyTrack.Album.Images),
			s.Key(),
			model.DefaultMarket,
			spotifyTrack.ExternalURLs["spotify"])
	}

	return res, res != nil, nil
}

func (s *spotifyStreamingService) GetPlaylistById(id string) (*model.Playlist, bool, error) {
	go s.metricsRecorder.CountSpotifyRequest()

	spotifyPlaylist, err := s.client.GetPlaylist(spotify.ID(id))
	if err != nil {
		return nil, false, err
	}

	res := model.NewPlaylist(
		spotifyPlaylist.ID.String(),
		spotifyPlaylist.Name,
		imageURL(spotifyPlaylist.Images),
		s.Key(),
		spotifyPlaylist.ExternalURLs["spotify"])

	return res, true, nil
}

func (s *spotifyStreamingService) GetPlaylistTracksById(id string) ([]*model.Track, bool, error) {

	var spotifyTracks []spotify.FullTrack

	// Todo: Cache these results, big playlists will make _a lot_ of requests
	offset := 0
	for {
		go s.metricsRecorder.CountSpotifyRequest()
		playlistTracks, err := s.client.GetPlaylistTracksOpt(spotify.ID(id), &spotify.Options{Offset: &offset}, "")
		if err != nil {
			return nil, false, err
		}

		for _, playlistTrack := range playlistTracks.Tracks {
			spotifyTracks = append(spotifyTracks, playlistTrack.Track)
		}

		if len(spotifyTracks) == playlistTracks.Total {
			break
		}

		offset = len(spotifyTracks)
	}

	var tracks []*model.Track
	for _, spotifyTrack := range spotifyTracks {
		track := model.NewTrack(
			spotifyTrack.ExternalIDs["isrc"],
			spotifyTrack.Name,
			artistName(spotifyTrack.Artists),
			spotifyTrack.Album.Name,
			imageURL(spotifyTrack.Album.Images),
			s.Key(),
			model.DefaultMarket,
			spotifyTrack.ExternalURLs["spotify"])
		tracks = append(tracks, track)
	}

	return tracks, true, nil
}

func (s *spotifyStreamingService) GetFromLink(link string) (model.Type, interface{}, error) {

	// example: https://open.spotify.com/track/4cOdK2wGLETKBW3PvgPWqT?si=10587ef152a8493f
	// format: 	https://open.spotify.com/<artist|album|track>/<id>?si=<user specific token that i don't care about>
	// Todo: How we gonna get the region?

	matches := findStringSubmatchMap(s.shareLinkPattern, link)

	// region := matches["region"]
	typ := matches["type"]
	id := spotify.ID(matches["id"])

	switch typ {
	case "artist":
		go s.metricsRecorder.CountSpotifyRequest()

		foundArtist, err := s.client.GetArtist(id)
		if err != nil {
			return model.UnknownType, false, err
		}

		artist := model.NewArtist(
			foundArtist.Name,
			imageURL(foundArtist.Images),
			s.Key(),
			model.DefaultMarket,
			foundArtist.ExternalURLs["spotify"])

		return model.ArtistType, artist, nil

	case "album":
		go s.metricsRecorder.CountSpotifyRequest()

		foundAlbum, err := s.client.GetAlbum(id)
		if err != nil {
			return model.UnknownType, nil, err
		}

		album := model.NewAlbum(
			foundAlbum.Name,
			artistName(foundAlbum.Artists),
			imageURL(foundAlbum.Images),
			s.Key(),
			model.DefaultMarket,
			foundAlbum.ExternalURLs["spotify"])

		return model.AlbumType, album, nil

	case "track":
		go s.metricsRecorder.CountSpotifyRequest()

		foundTrack, err := s.client.GetTrack(id)
		if err != nil {
			return model.UnknownType, nil, err
		}

		track := model.NewTrack(
			foundTrack.ExternalIDs["isrc"],
			foundTrack.Name,
			artistName(foundTrack.Artists),
			foundTrack.Album.Name,
			imageURL(foundTrack.Album.Images),
			s.Key(),
			model.DefaultMarket,
			foundTrack.ExternalURLs["spotify"])

		return model.TrackType, track, nil

	case "playlist":
		go s.metricsRecorder.CountSpotifyRequest()

		foundPlaylist, err := s.client.GetPlaylist(id)
		if err != nil {
			return model.UnknownType, nil, err
		}

		playlist := model.NewPlaylist(
			foundPlaylist.ID.String(),
			foundPlaylist.Name,
			imageURL(foundPlaylist.Images),
			s.Key(),
			foundPlaylist.ExternalURLs["spotify"])

		return model.PlaylistType, playlist, nil

	default:
		return model.UnknownType, nil, fmt.Errorf("unknown type %s", typ)
	}
}

func (s *spotifyStreamingService) CleanLink(link string) string {

	match := s.shareLinkPattern.FindStringIndex(link)
	if len(match) > 0 {
		return link[match[0]:match[1]]
	}

	return link
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

func imageURL(imgs []spotify.Image) string {
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
