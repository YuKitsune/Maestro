package appleMusic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

const baseUrl = "https://api.music.apple.com"

type appleMusicStreamingService struct {
	c                *http.Client
	shareLinkPattern *regexp.Regexp
}

func NewAppleMusicStreamingService(token string, shareLinkPattern string) streamingService.StreamingService {
	shareLinkPatternRegex := regexp.MustCompile(shareLinkPattern)
	return &appleMusicStreamingService{
		streamingService.NewClientWithBearerAuth(token),
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

	storefront := artist.GetMarket()
	term := strings.ReplaceAll(artist.Name, " ", "+")
	apiUrl := fmt.Sprintf("%s/v1/catalog/%s/search?term=%s&types=artists", baseUrl, storefront, term)

	httpRes, err := s.c.Get(apiUrl)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var apiRes *SearchResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return nil, err
	}

	if len(apiRes.Results.Artists.Data) == 0 {
		return nil, nil
	}

	appleMusicArtist := apiRes.Results.Artists.Data[0]
	res := model.NewArtist(
		appleMusicArtist.Attributes.Name,
		"", // Todo: artwork link
		s.Name(),
		storefront,
		appleMusicArtist.Attributes.Url)
	return res, nil
}

func (s *appleMusicStreamingService) SearchAlbum(album *model.Album) (*model.Album, error) {

	storefront := album.GetMarket()
	term := strings.ReplaceAll(album.Name, " ", "+")
	apiUrl := fmt.Sprintf("%s/v1/catalog/%s/search?term=%s&types=albums", baseUrl, storefront, term)

	httpRes, err := s.c.Get(apiUrl)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var apiRes *SearchResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return nil, err
	}

	if len(apiRes.Results.Albums.Data) == 0 {
		return nil, nil
	}

	appleMusicAlbum := apiRes.Results.Albums.Data[0]

	// Query relationships for artist names
	artistNames, err := s.getArtistNames(appleMusicAlbum.Relationships)
	if err != nil {
		return nil, err
	}

	res := model.NewAlbum(
		normalizeAlbumName(appleMusicAlbum),
		artistNames,
		appleMusicAlbum.Attributes.Artwork.Url,
		s.Name(),
		storefront,
		appleMusicAlbum.Attributes.Url)

	return res, nil
}

func (s *appleMusicStreamingService) SearchSong(song *model.Track) (*model.Track, error) {

	storefront := song.GetMarket()
	rawTerm := fmt.Sprintf("%s - %s", strings.Join(song.ArtistNames, ", "), song.Name)
	term := strings.ReplaceAll(rawTerm, " ", "+")
	apiUrl := fmt.Sprintf("%s/v1/catalog/%s/search?term=%s&types=songs", baseUrl, storefront, term)

	httpRes, err := s.c.Get(apiUrl)
	defer httpRes.Body.Close()

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var apiRes *SearchResponse
	err = json.Unmarshal(resBytes, &apiRes)
	if err != nil {
		return nil, err
	}

	if len(apiRes.Results.Songs.Data) == 0 {
		return nil, nil
	}

	// Todo: Not enough info in initial query, need to filter down here...
	appleMusicSong, err := s.filterSongs(apiRes.Results.Songs, song)
	if err != nil {
		return nil, err
	}

	// Query relationships for artist names
	artistNames, err := s.getArtistNames(appleMusicSong.Relationships)
	if err != nil {
		return nil, err
	}

	// Query relationships for parent album
	// Todo: What if there are many?
	albumName, err := s.getAlbumName(appleMusicSong.Relationships, song.AlbumName)
	if err != nil {
		return nil, err
	}

	res := model.NewTrack(
		appleMusicSong.Attributes.Name,
		artistNames,
		albumName,
		s.Name(),
		storefront,
		appleMusicSong.Attributes.Url)

	return res, nil
}

func (s *appleMusicStreamingService) SearchFromLink(link string) (model.Thing, error) {

	// example: https://music.apple.com/au/album/surrender/1585865534
	// format: 	https://music.apple.com/<storefront>/<artist|album|song>/<name>/<id>
	// name is irrelevant here, we only need the storefront and id

	matches := streamingService.FindStringSubmatchMap(s.shareLinkPattern, link)

	storefront := model.Market(matches["storefront"])
	typ := matches["type"]
	id := matches["id"]

	var apiUrl string
	var unmarshalFunc func(rb []byte) (model.Thing, error)

	switch typ {
	case "artist":
		apiUrl = fmt.Sprintf("%s/v1/catalog/%s/artists/%s", baseUrl, storefront, id)
		unmarshalFunc = func(rb []byte) (model.Thing, error) {
			var apiRes *ArtistsResult
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			if len(apiRes.Data) == 0 {
				return nil, fmt.Errorf("no %s found in region %s with id %s", typ, storefront, id)
			}

			foundArtist := apiRes.Data[0]
			artist := model.NewArtist(
				foundArtist.Attributes.Name,
				"", // Todo: artwork link
				s.Name(),
				storefront,
				foundArtist.Attributes.Url)

			return artist, nil
		}

		break

	case "album":
		apiUrl = fmt.Sprintf("%s/v1/catalog/%s/albums/%s", baseUrl, storefront, id)
		unmarshalFunc = func(rb []byte) (model.Thing, error) {
			var apiRes *AlbumResult
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			foundAlbum := apiRes.Data[0]

			albumName := normalizeAlbumName(foundAlbum)

			// Query relationships for artist names
			artistNames, err := s.getArtistNames(foundAlbum.Relationships)
			if err != nil {
				return nil, err
			}

			album := model.NewAlbum(
				albumName,
				artistNames,
				foundAlbum.Attributes.Artwork.Url,
				s.Name(),
				storefront,
				foundAlbum.Attributes.Url)

			return album, nil
		}
		break

	case "song":
		apiUrl = fmt.Sprintf("%s/v1/catalog/%s/songs/%s", baseUrl, storefront, id)
		unmarshalFunc = func(rb []byte) (model.Thing, error) {
			var apiRes *SongResult
			err := json.Unmarshal(rb, &apiRes)
			if err != nil {
				return nil, err
			}

			foundSong := apiRes.Data[0]

			// Query relationships for artist names
			artistNames, err := s.getArtistNames(foundSong.Relationships)
			if err != nil {
				return nil, err
			}

			// Query relationships for parent album
			// Todo: What if there are many?
			albumName, err := s.getAlbumName(foundSong.Relationships)
			if err != nil {
				return nil, err
			}

			song := model.NewTrack(
				foundSong.Attributes.Name,
				artistNames,
				albumName,
				s.Name(),
				storefront,
				foundSong.Attributes.Url)

			return song, nil
		}
		break

	default:
		return nil, fmt.Errorf("unknown type %s", typ)
	}

	httpRes, err := s.c.Get(apiUrl)
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no %s found in region %s with id %s", typ, storefront, id)
	}

	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	res, err := unmarshalFunc(resBytes)

	return res, err
}

func normalizeAlbumName(album *Album) string {
	suffix := " - Single"
	suffixLen := len(suffix)

	name := album.Attributes.Name
	nameLen := len(name)

	if album.Attributes.IsSingle && strings.HasSuffix(name, suffix) {
		normName := name[0:nameLen - suffixLen]
		return normName
	}

	return album.Attributes.Name
}

func (s *appleMusicStreamingService) getArtistNames(rel *Relationships) ([]string, error) {
	var names []string
	for _, data := range rel.Artists.Data {
		if data.Type != "artists" {
			continue
		}

		url := fmt.Sprintf("%s/%s", baseUrl, data.Href)
		httpRes, err := s.c.Get(url)
		defer httpRes.Body.Close()
		if err != nil {
			return nil, err
		}

		resBytes, err := ioutil.ReadAll(httpRes.Body)
		if err != nil {
			return nil, err
		}

		var res *ArtistsResult
		err = json.Unmarshal(resBytes, &res)
		if err != nil {
			return nil, err
		}

		names = append(names, res.Data[0].Attributes.Name)
	}

	return names, nil
}

func (s *appleMusicStreamingService) getAlbumName(rel *Relationships, targetAlbumNames ...string) (string, error) {
	for _, data := range rel.Albums.Data {
		if data.Type != "albums" {
			continue
		}

		url := fmt.Sprintf("%s/%s", baseUrl, data.Href)
		httpRes, err := s.c.Get(url)
		defer httpRes.Body.Close()
		if err != nil {
			return "", err
		}

		resBytes, err := ioutil.ReadAll(httpRes.Body)
		if err != nil {
			return "", err
		}

		var res *AlbumResult
		err = json.Unmarshal(resBytes, &res)
		if err != nil {
			return "", err
		}

		album := res.Data[0]

		// Yuck...
		normName := normalizeAlbumName(album)
		if len(targetAlbumNames) > 0 {
			for _, name := range targetAlbumNames {
				if len(name) > 0 && normName == name {
					return normName, nil
				}
			}
		} else {
			return normName, nil
		}
	}

	return "", nil
}

func (s *appleMusicStreamingService) filterSongs(songResult *SongResult, targetTrack *model.Track) (*Song, error) {

	for _, songSearchResult := range songResult.Data {

		httpRes, err := s.c.Get(fmt.Sprintf("%s/v1/catalog/%s/songs/%s", baseUrl, targetTrack.Market, songSearchResult.Id))
		if err != nil {
			return nil, err
		}
		defer httpRes.Body.Close()

		resBytes, err := ioutil.ReadAll(httpRes.Body)
		if err != nil {
			return nil, err
		}

		var sr *SongResult
		err = json.Unmarshal(resBytes, &sr)
		if err != nil {
			return nil, err
		}

		song := sr.Data[0]

		namesMatch := song.Attributes.Name == targetTrack.Name

		songArtists, err := s.getArtistNames(song.Relationships)
		if err != nil {
			return nil, err
		}

		targetArtists := targetTrack.ArtistNames
		artistsMatch := true
		sort.Strings(songArtists)
		for _, artist := range targetArtists {
			if sort.SearchStrings(songArtists, artist) == len(songArtists) {
				artistsMatch = false
			}
		}

		albumsMatch := true
		if len(targetTrack.AlbumName) > 0 {
			songAlbum, err := s.getAlbumName(song.Relationships, targetTrack.AlbumName)
			if err != nil {
				return nil, err
			}

			if len(songAlbum) > 0 {
				albumsMatch = songAlbum == targetTrack.AlbumName
			}
		}

		if namesMatch && artistsMatch && albumsMatch {
			return song, nil
		}
	}

	return nil, nil
}