package appleMusic

import (
	"encoding/json"
	"fmt"
	"github.com/noppefoxwolf/amg/applemusic"
	"io/ioutil"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"regexp"
	"sort"
	"strings"
)

type appleMusicStreamingService struct {
	client *applemusic.Client
	shareLinkPattern *regexp.Regexp
}

func NewAppleMusicStreamingService(token string, shareLinkPattern string) streamingService.StreamingService {
	shareLinkPatternRegex := regexp.MustCompile(shareLinkPattern)

	c := streamingService.NewClientWithBearerAuth(token)
	amc := applemusic.NewClient(c)

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

	p := &applemusic.SearchForCatalogResourcesParams{
		Storefront: artist.Market.String(),
		Term:       artist.Name,
		Types: []string{"artists"},
	}

	amRes, _, err := s.client.Search.SearchForCatalogResources(p)
	if err != nil {
		return nil, err
	}

	var foundArtists []*applemusic.Artist
	for _, data := range amRes.Data {
		foundArtists = append(foundArtists, data.Artist)
	}

	// Todo: Narrow down results
	foundArtist := foundArtists[0]

	return s.newArtist(foundArtist, artist.Market)
}

func (s *appleMusicStreamingService) SearchAlbum(album *model.Album) (*model.Album, error) {

	p := &applemusic.SearchForCatalogResourcesParams{
		Storefront: album.Market.String(),
		Term:       fmt.Sprintf("%s %s", strings.Join(album.ArtistNames, ", "), album.Name),
		Types: []string{"albums"},
	}

	amRes, _, err := s.client.Search.SearchForCatalogResources(p)
	if err != nil {
		return nil, err
	}

	var foundAlbums []*applemusic.Album
	for _, data := range amRes.Data {
		foundAlbums = append(foundAlbums, data.Album)
	}

	// Todo: Narrow down results
	foundAlbum := foundAlbums[0]

	return s.newAlbum(foundAlbum, album.Market)
}

func (s *appleMusicStreamingService) SearchSong(song *model.Track) (*model.Track, error) {

	p := &applemusic.SearchForCatalogResourcesParams{
		Storefront: song.Market.String(),
		Term:       fmt.Sprintf("%s - %s", strings.Join(song.ArtistNames, ", "), song.Name),
		Types: []string{"songs"},
	}

	amRes, _, err := s.client.Search.SearchForCatalogResources(p)
	if err != nil {
		return nil, err
	}

	var foundSongs []*applemusic.Song
	for _, data := range amRes.Data {
		foundSongs = append(foundSongs, data.Song)
	}

	// Todo: Narrow down results
	foundSong := foundSongs[0]

	return s.newTrack(foundSong, song.Market)
}

func (s *appleMusicStreamingService) SearchFromLink(link string) (model.Thing, error) {

	// example: https://music.apple.com/au/album/surrender/1585865534
	// format: 	https://music.apple.com/<storefront>/<artist|album|song>/<name>/<id>
	// name is irrelevant here, we only need the storefront and id

	matches := streamingService.FindStringSubmatchMap(s.shareLinkPattern, link)

	storefront := model.Market(matches["storefront"])
	typ := matches["type"]
	id := matches["id"]

	switch typ {
	case "artist":
		amRes, _, err := s.client.Artists.GetACatalogArtist(&applemusic.GetACatalogArtistParams{
			Id:         id,
			Storefront: storefront.String(),
		})
		if err != nil {
			return nil, err
		}

		var foundArtists []*applemusic.Artist
		for _, data := range amRes.Data {
			foundArtists = append(foundArtists, data.Artist)
		}

		// Todo: Narrow down results
		foundArtist := foundArtists[0]

		return s.newArtist(foundArtist, storefront)

	case "album":
		amRes, _, err := s.client.Albums.GetACatalogAlbum(&applemusic.GetACatalogAlbumParams{
			Id:         id,
			Storefront: storefront.String(),
			Include: []string{"artists"},
		})
		if err != nil {
			return nil, err
		}

		var foundAlbums []*applemusic.Album
		for _, data := range amRes.Data {
			foundAlbums = append(foundAlbums, data.Album)
		}

		// Todo: Narrow down results
		foundAlbum := foundAlbums[0]

		return s.newAlbum(foundAlbum, storefront)

	case "track":
		amRes, _, err := s.client.Songs.GetACatalogSong(&applemusic.GetACatalogSongParams{
			Id:         id,
			Storefront: storefront.String(),
			Include: []string{"artists", "albums"},
		})
		if err != nil {
			return nil, err
		}

		var foundSongs []*applemusic.Song
		for _, data := range amRes.Data {
			foundSongs = append(foundSongs, data.Song)
		}

		// Todo: Narrow down results
		foundSong := foundSongs[0]

		return s.newTrack(foundSong, storefront)

	default:
		return nil, fmt.Errorf("unknown type %s", typ)
	}
}

func (s *appleMusicStreamingService) newArtist(artist *applemusic.Artist, market model.Market) (*model.Artist, error) {

	newArtist := model.NewArtist(
		artist.Attributes.Name,
		"", // Todo: artwork link
		s.Name(),
		market,
		artist.Attributes.Url)

	return newArtist, nil
}

func (s *appleMusicStreamingService) newAlbum(album *applemusic.Album, market model.Market) (*model.Album, error) {

	// Query relationships for artist names
	artistNames, err := s.getAlbumArtistNames(album)
	if err != nil {
		return nil, err
	}

	newAlbum := model.NewAlbum(
		album.Attributes.Name,
		artistNames,
		album.Attributes.Artwork.Url,
		s.Name(),
		market,
		album.Attributes.Url)

	return newAlbum, nil
}

func (s *appleMusicStreamingService) newTrack(song *applemusic.Song, market model.Market) (*model.Track, error) {

	// Query relationships for artist names
	artistNames, err := s.getSongArtistNames(song)
	if err != nil {
		return nil, err
	}

	// Query relationships for album
	// Todo: What if there are many?
	albumName, err := s.getSongAlbumName(song)
	if err != nil {
		return nil, err
	}

	track := model.NewTrack(
		song.Attributes.Name,
		artistNames,
		albumName,
		s.Name(),
		market,
		song.Attributes.Url)

	return track, nil
}

//func normalizeAlbumName(album *Album) string {
//	suffix := " - Single"
//	suffixLen := len(suffix)
//
//	name := album.Attributes.Name
//	nameLen := len(name)
//
//	if album.Attributes.IsSingle && strings.HasSuffix(name, suffix) {
//		normName := name[0:nameLen - suffixLen]
//		return normName
//	}
//
//	return album.Attributes.Name
//}
//
//func (s *appleMusicStreamingService) getArtistNames(rel *applemusic.ArtistRelationships, fallback string) ([]string, error) {
//	var names []string
//
//	if rel == nil || rel.Artists.Data == nil || len(rel.Artists.Data) == 0 {
//		return []string {fallback}, nil
//	}
//
//	for _, data := range rel.Artists.Data {
//		if data.Type != "artists" {
//			continue
//		}
//
//		url := fmt.Sprintf("%s/%s", baseUrl, data.Href)
//		httpRes, err := s.c.Get(url)
//		defer httpRes.Body.Close()
//		if err != nil {
//			return nil, err
//		}
//
//		resBytes, err := ioutil.ReadAll(httpRes.Body)
//		if err != nil {
//			return nil, err
//		}
//
//		var res *ArtistsResult
//		err = json.Unmarshal(resBytes, &res)
//		if err != nil {
//			return nil, err
//		}
//
//		names = append(names, res.Data[0].Attributes.Name)
//	}
//
//	return names, nil
//}
//
//func (s *appleMusicStreamingService) getAlbumName(rel *Relationships, targetAlbumNames ...string) (string, error) {
//	for _, data := range rel.Albums.Data {
//		if data.Type != "albums" {
//			continue
//		}
//
//		url := fmt.Sprintf("%s/%s", baseUrl, data.Href)
//		httpRes, err := s.c.Get(url)
//		defer httpRes.Body.Close()
//		if err != nil {
//			return "", err
//		}
//
//		resBytes, err := ioutil.ReadAll(httpRes.Body)
//		if err != nil {
//			return "", err
//		}
//
//		var res *AlbumResult
//		err = json.Unmarshal(resBytes, &res)
//		if err != nil {
//			return "", err
//		}
//
//		album := res.Data[0]
//
//		// Yuck...
//		normName := album.Attributes.Name
//		if len(targetAlbumNames) > 0 {
//			for _, name := range targetAlbumNames {
//				if len(name) > 0 && normName == name {
//					return normName, nil
//				}
//			}
//		} else {
//			return normName, nil
//		}
//	}
//
//	return "", nil
//}
//
//func (s *appleMusicStreamingService) filterSongs(songResult *SongResult, targetTrack *model.Track) (*Song, error) {
//
//	for _, songSearchResult := range songResult.Data {
//
//		httpRes, err := s.c.Get(fmt.Sprintf("%s/v1/catalog/%s/songs/%s", baseUrl, targetTrack.Market, songSearchResult.Id))
//		if err != nil {
//			return nil, err
//		}
//		defer httpRes.Body.Close()
//
//		resBytes, err := ioutil.ReadAll(httpRes.Body)
//		if err != nil {
//			return nil, err
//		}
//
//		var sr *SongResult
//		err = json.Unmarshal(resBytes, &sr)
//		if err != nil {
//			return nil, err
//		}
//
//		song := sr.Data[0]
//
//		namesMatch := song.Attributes.Name == targetTrack.Name
//
//		songArtists, err := s.getArtistNames(song.Relationships, song.Attributes.ArtistName)
//		if err != nil {
//			return nil, err
//		}
//
//		targetArtists := targetTrack.ArtistNames
//		artistsMatch := true
//		sort.Strings(songArtists)
//		for _, artist := range targetArtists {
//			if sort.SearchStrings(songArtists, artist) == len(songArtists) {
//				artistsMatch = false
//			}
//		}
//
//		albumsMatch := true
//		if len(targetTrack.AlbumName) > 0 {
//			songAlbum, err := s.getAlbumName(song.Relationships, targetTrack.AlbumName)
//			if err != nil {
//				return nil, err
//			}
//
//			if len(songAlbum) > 0 {
//				albumsMatch = songAlbum == targetTrack.AlbumName
//			}
//		}
//
//		if namesMatch && artistsMatch && albumsMatch {
//			return song, nil
//		}
//	}
//
//	return nil, nil
//}