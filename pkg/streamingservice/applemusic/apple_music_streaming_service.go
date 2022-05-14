package applemusic

import (
	"fmt"
	"github.com/yukitsune/maestro/pkg/metrics"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingservice"
	"regexp"
	"strings"
)

type appleMusicStreamingService struct {
	client           *client
	shareLinkPattern *regexp.Regexp
	metricsRecorder  metrics.Recorder
}

func NewAppleMusicStreamingService(cfg *Config, mr metrics.Recorder) streamingservice.StreamingService {
	shareLinkPatternRegex := regexp.MustCompile("(https?:\\/\\/)?music\\.apple\\.com\\/(?P<storefront>[A-Za-z0-9]+)\\/(?P<type>[A-Za-z]+)\\/(?:.+\\/)(?P<id>[0-9]+)(?:\\?i=(?P<song_id>[0-9]+))?")

	amc := NewAppleMusicClient(cfg.Token)

	return &appleMusicStreamingService{
		amc,
		shareLinkPatternRegex,
		mr,
	}
}

func (s *appleMusicStreamingService) Key() model.StreamingServiceKey {
	return Key
}

func (s *appleMusicStreamingService) LinkBelongsToService(link string) bool {
	return s.shareLinkPattern.MatchString(link)
}

func (s *appleMusicStreamingService) SearchArtist(artist *model.Artist) (*model.Artist, bool, error) {

	go s.metricsRecorder.CountAppleMusicRequest()

	searchRes, err := s.client.SearchArtist(artist.Name, artist.GetMarket())
	if err != nil {
		return nil, false, err
	}

	if len(searchRes) == 0 {
		return nil, false, nil
	}

	// Todo: Narrow down results
	foundArtist := searchRes[0]

	artistRes, err := s.newArtist(&foundArtist, artist.Market)
	return artistRes, true, err
}

func (s *appleMusicStreamingService) SearchAlbum(album *model.Album) (*model.Album, bool, error) {

	go s.metricsRecorder.CountAppleMusicRequest()

	term := fmt.Sprintf("%s %s", strings.Join(album.ArtistNames, " "), album.Name)
	searchRes, err := s.client.SearchAlbum(term, album.GetMarket())
	if err != nil {
		return nil, false, err
	}

	if len(searchRes) == 0 {
		return nil, false, nil
	}

	// Todo: Narrow down results
	foundAlbum := searchRes[0]

	// Load the album directly so we get the relationships
	fullAlbum, err := s.client.GetAlbum(foundAlbum.ID, album.Market)
	if err != nil {
		return nil, false, err
	}

	resAlbum, err := s.newAlbum(fullAlbum, album.Market)
	return resAlbum, true, err
}

func (s *appleMusicStreamingService) GetTrackByIsrc(isrc string) (*model.Track, bool, error) {

	songsRes, err := s.client.GetSongByIsrc(isrc, model.DefaultMarket)
	if err != nil {
		return nil, false, err
	}

	if len(songsRes) == 0 {
		return nil, false, nil
	}

	// Todo: Narrow down results
	foundSong := songsRes[0]

	track, err := s.newTrack(&foundSong, model.DefaultMarket)
	if err != nil {
		return nil, false, err
	}

	return track, true, nil
}

func (s *appleMusicStreamingService) SearchTrack(song *model.Track) (*model.Track, bool, error) {

	go s.metricsRecorder.CountAppleMusicRequest()

	var searchRes []Song
	var err error
	if len(song.Isrc) > 0 {
		searchRes, err = s.client.GetSongByIsrc(song.Isrc, song.Market)
		if err != nil {
			return nil, false, err
		}
	} else {
		term := fmt.Sprintf("%s %s", strings.Join(song.ArtistNames, " "), song.Name)
		searchRes, err = s.client.SearchSong(term, song.GetMarket())
		if err != nil {
			return nil, false, err
		}
	}

	if len(searchRes) == 0 {
		return nil, false, nil
	}

	// Todo: Narrow down results
	foundSong := searchRes[0]

	// Load the song directly so we get the relationships
	fullSong, err := s.client.GetSong(foundSong.ID, song.Market)
	if err != nil {
		return nil, false, err
	}

	resTrack, err := s.newTrack(fullSong, song.Market)
	return resTrack, true, err
}

func (s *appleMusicStreamingService) GetPlaylistById(id string) (*model.Playlist, bool, error) {
	panic("not implemented")
}

func (s *appleMusicStreamingService) GetPlaylistTracksById(id string) ([]*model.Track, bool, error) {
	panic("not implemented")
}

func (s *appleMusicStreamingService) GetFromLink(link string) (model.Type, interface{}, error) {

	// example: https://music.apple.com/au/album/surrender/1585865534?i=123123123
	// format: 	https://music.apple.com/<storefront>/<artist|album>/<name>/<album-id/artist-id>?i=<song-id>
	// name is irrelevant here, we only need the storefront, type, and ids

	matches := streamingservice.FindStringSubmatchMap(s.shareLinkPattern, link)

	storefront := model.Market(matches["storefront"])
	typStr := matches["type"]
	id := matches["id"]
	songID := matches["song_id"]

	// Tracks/Songs don't have their own links,
	// they're links to albums with a songID attached
	if typStr == "album" && len(songID) > 0 {
		typStr = "song"
		id = songID
	}

	switch typStr {
	case "artist":
		go s.metricsRecorder.CountAppleMusicRequest()
		res, err := s.client.GetArtist(id, storefront)
		if err != nil {
			return model.UnknownType, nil, err
		}

		artist, err := s.newArtist(res, storefront)
		return model.ArtistType, artist, err

	case "album":
		go s.metricsRecorder.CountAppleMusicRequest()
		res, err := s.client.GetAlbum(id, storefront)
		if err != nil {
			return model.UnknownType, nil, err
		}

		album, err := s.newAlbum(res, storefront)
		return model.AlbumType, album, err

	case "song":
		go s.metricsRecorder.CountAppleMusicRequest()
		res, err := s.client.GetSong(id, storefront)
		if err != nil {
			return model.UnknownType, nil, err
		}

		track, err := s.newTrack(res, storefront)
		return model.TrackType, track, err

	case "playlist":
		panic("not implemented")

	default:
		return model.UnknownType, nil, fmt.Errorf("unknown type %s", typStr)
	}
}

func (s *appleMusicStreamingService) CleanLink(link string) string {

	match := s.shareLinkPattern.FindStringIndex(link)
	if len(match) > 0 {
		return link[match[0]:match[1]]
	}

	return link
}

func (s *appleMusicStreamingService) newArtist(artist *Artist, market model.Market) (*model.Artist, error) {

	newArtist := model.NewArtist(
		artist.Attributes.Name,
		"",
		s.Key(),
		market,
		artist.Attributes.URL)

	return newArtist, nil
}

func (s *appleMusicStreamingService) newAlbum(album *Album, market model.Market) (*model.Album, error) {

	// Clean up album name
	// Todo: Revisit
	albumName := album.Attributes.Name
	if album.Attributes.IsSingle {
		singleRegex := regexp.MustCompile("\\s-\\sSingle$")
		indexes := singleRegex.FindStringIndex(albumName)
		if len(indexes) > 0 {
			albumName = albumName[0:indexes[0]]
		}
	}

	// Query relationships for artist names
	artistNames, err := s.getAlbumArtistNames(album, market)
	if err != nil {
		return nil, err
	}

	newAlbum := model.NewAlbum(
		albumName,
		artistNames,
		getArtworkURL(&album.Attributes.Artwork),
		s.Key(),
		market,
		album.Attributes.URL)

	return newAlbum, nil
}

func (s *appleMusicStreamingService) newTrack(song *Song, market model.Market) (*model.Track, error) {

	// Query relationships for artist names
	artistNames, err := s.getSongArtistNames(song, market)
	if err != nil {
		return nil, err
	}

	// Query relationships for album artwork
	artworkLink, err := s.getSongArtwork(song, market)
	if err != nil {
		return nil, err
	}

	track := model.NewTrack(
		song.Attributes.Isrc,
		song.Attributes.Name,
		artistNames,
		song.Attributes.AlbumName,
		artworkLink,
		s.Key(),
		market,
		song.Attributes.URL)

	return track, nil
}

func getArtworkURL(art *Artwork) string {
	url := art.URL
	url = strings.ReplaceAll(url, "{w}", fmt.Sprintf("%d", art.Width))
	url = strings.ReplaceAll(url, "{h}", fmt.Sprintf("%d", art.Height))

	return url
}

func (s *appleMusicStreamingService) getAlbumArtistNames(album *Album, market model.Market) ([]string, error) {
	var names []string

	for _, data := range album.Relationships.Artists.Data {

		go s.metricsRecorder.CountAppleMusicRequest()

		artist, err := s.client.GetArtist(data.ID, market)
		if err != nil {
			return names, nil
		}

		names = append(names, artist.Attributes.Name)
	}

	return names, nil
}

func (s *appleMusicStreamingService) getSongArtistNames(song *Song, market model.Market) ([]string, error) {
	var names []string

	for _, data := range song.Relationships.Artists.Data {

		go s.metricsRecorder.CountAppleMusicRequest()

		artist, err := s.client.GetArtist(data.ID, market)
		if err != nil {
			return names, nil
		}

		names = append(names, artist.Attributes.Name)
	}

	return names, nil
}

func (s *appleMusicStreamingService) getSongArtwork(song *Song, market model.Market) (string, error) {

	var artworkLink string
	if len(song.Relationships.Albums.Data) > 0 {

		data := song.Relationships.Albums.Data[0]

		go s.metricsRecorder.CountAppleMusicRequest()

		album, err := s.client.GetAlbum(data.ID, market)
		if err != nil {
			return artworkLink, err
		}

		artworkLink = getArtworkURL(&album.Attributes.Artwork)
	}

	return artworkLink, nil
}
