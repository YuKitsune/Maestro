package appleMusic

import (
	"fmt"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"regexp"
	"strings"
)

type appleMusicStreamingService struct {
	client           *AppleMusicClient
	shareLinkPattern *regexp.Regexp
}

func NewAppleMusicStreamingService(cfg *Config) streamingService.StreamingService {
	shareLinkPatternRegex := regexp.MustCompile("(?:https:\\/\\/music\\.apple\\.com\\/)(?P<storefront>[A-Za-z0-9]+)\\/(?P<type>[A-Za-z]+)\\/(?:.+\\/)(?P<id>[0-9]+)(?:\\?i=(?P<song_id>[0-9]+))?")

	amc := NewAppleMusicClient(cfg.Token)

	return &appleMusicStreamingService{
		amc,
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

	searchRes, err := s.client.SearchArtist(artist.Name, artist.GetMarket())
	if err != nil {
		return nil, err
	}

	// Todo: Narrow down results
	foundArtist := searchRes[0]

	return s.newArtist(&foundArtist, artist.Market)
}

func (s *appleMusicStreamingService) SearchAlbum(album *model.Album) (*model.Album, error) {

	term := fmt.Sprintf("%s %s", strings.Join(album.ArtistNames, " "), album.Name)
	searchRes, err := s.client.SearchAlbum(term, album.GetMarket())
	if err != nil {
		return nil, err
	}

	// Todo: Narrow down results
	foundAlbum := searchRes[0]

	// Load the album directly so we get the relationships
	fullAlbum, err := s.client.GetAlbum(foundAlbum.Id, album.Market)
	if err != nil {
		return nil, err
	}

	return s.newAlbum(fullAlbum, album.Market)
}

func (s *appleMusicStreamingService) SearchSong(song *model.Track) (*model.Track, error) {

	term := fmt.Sprintf("%s %s", strings.Join(song.ArtistNames, " "), song.Name)
	searchRes, err := s.client.SearchSong(term, song.GetMarket())
	if err != nil {
		return nil, err
	}

	// Todo: Narrow down results
	foundSong := searchRes[0]

	// Load the song directly so we get the relationships
	fullSong, err := s.client.GetSong(foundSong.Id, song.Market)
	if err != nil {
		return nil, err
	}

	return s.newTrack(fullSong, song.Market)
}

func (s *appleMusicStreamingService) SearchFromLink(link string) (model.Thing, error) {

	// example: https://music.apple.com/au/album/surrender/1585865534?i=123123123
	// format: 	https://music.apple.com/<storefront>/<artist|album>/<name>/<album-id/artist-id>?i=<song-id>
	// name is irrelevant here, we only need the storefront, type, and ids

	matches := streamingService.FindStringSubmatchMap(s.shareLinkPattern, link)

	storefront := model.Market(matches["storefront"])
	typ := matches["type"]
	id := matches["id"]
	songId := matches["song_id"]

	// Hack but it works
	if typ == "album" && len(songId) > 0 {
		typ = "song"
		id = songId
	}

	switch typ {
	case "artist":
		res, err := s.client.GetArtist(id, storefront)
		if err != nil {
			return nil, err
		}

		return s.newArtist(res, storefront)

	case "album":
		res, err := s.client.GetAlbum(id, storefront)
		if err != nil {
			return nil, err
		}

		return s.newAlbum(res, storefront)

	case "song":
		res, err := s.client.GetSong(id, storefront)
		if err != nil {
			return nil, err
		}

		return s.newTrack(res, storefront)

	default:
		return nil, fmt.Errorf("unknown type %s", typ)
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
		s.Name(),
		market,
		artist.Attributes.Url)

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
		getArtworkUrl(&album.Attributes.Artwork),
		s.Name(),
		market,
		album.Attributes.Url)

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
		song.Attributes.Name,
		artistNames,
		song.Attributes.AlbumName,
		artworkLink,
		s.Name(),
		market,
		song.Attributes.Url)

	return track, nil
}

func getArtworkUrl(art *Artwork) string {
	url := art.Url
	url = strings.ReplaceAll(url, "{w}", fmt.Sprintf("%d", art.Width))
	url = strings.ReplaceAll(url, "{h}", fmt.Sprintf("%d", art.Height))

	return url
}

func (s *appleMusicStreamingService) getAlbumArtistNames(album *Album, market model.Market) ([]string, error) {
	var names []string

	for _, data := range album.Relationships.Artists.Data {
		artist, err := s.client.GetArtist(data.Id, market)
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
		artist, err := s.client.GetArtist(data.Id, market)
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

		album, err := s.client.GetAlbum(data.Id, market)
		if err != nil {
			return artworkLink, err
		}

		artworkLink = getArtworkUrl(&album.Attributes.Artwork)
	}

	return artworkLink, nil
}