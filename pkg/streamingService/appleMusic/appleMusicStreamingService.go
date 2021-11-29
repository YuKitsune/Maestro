package appleMusic

import (
	"fmt"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"regexp"
	"strings"
)

type appleMusicStreamingService struct {
	client *AppleMusicClient
	shareLinkPattern *regexp.Regexp
}

func NewAppleMusicStreamingService(token string, shareLinkPattern string) streamingService.StreamingService {
	shareLinkPatternRegex := regexp.MustCompile(shareLinkPattern)

	amc := NewAppleMusicClient(token)

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

	return s.newAlbum(&foundAlbum, album.Market)
}

func (s *appleMusicStreamingService) SearchSong(song *model.Track) (*model.Track, error) {

	term := fmt.Sprintf("%s %s", strings.Join(song.ArtistNames, " "), song.Name)
	searchRes, err := s.client.SearchSong(term, song.GetMarket())
	if err != nil {
		return nil, err
	}

	// Todo: Narrow down results
	foundSong := searchRes[0]

	return s.newTrack(&foundSong, song.Market)
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

func (s *appleMusicStreamingService) newArtist(artist *Artist, market model.Market) (*model.Artist, error) {

	newArtist := model.NewArtist(
		artist.Attributes.Name,
		"", // Todo: artwork link
		s.Name(),
		market,
		artist.Attributes.Url)

	return newArtist, nil
}

func (s *appleMusicStreamingService) newAlbum(album *Album, market model.Market) (*model.Album, error) {

	// Query relationships for artist names
	artistNames, err := s.getAlbumArtistNames(album, market)
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

func (s *appleMusicStreamingService) newTrack(song *Song, market model.Market) (*model.Track, error) {

	// Query relationships for artist names
	artistNames, err := s.getSongArtistNames(song, market)
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

// TODO: Figure out what we should do here
// 	is song.Attributes.AlbumName sufficent?
// 	do we need to check through the attributes?

func (s *appleMusicStreamingService) getSongAlbumName(song *Song) (string, error) {
	return song.Attributes.AlbumName, nil
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