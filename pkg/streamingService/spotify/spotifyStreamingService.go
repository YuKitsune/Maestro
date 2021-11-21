package spotify

import (
	"github.com/zmb3/spotify"
	"maestro/pkg/streamingService"
)

type spotifyStreamingService struct {
	client *spotify.Client
}

func NewSpotifyStreamingService(token string) streamingService.StreamingService {
	c := streamingService.NewClientWithBearerAuth(token)
	sc := spotify.NewClient(c)
	return &spotifyStreamingService{client: &sc}
}

func (s *spotifyStreamingService) Name() string {
	return "Spotify"
}

func (s *spotifyStreamingService) SearchArtist(name string) (res []streamingService.Artist, err error) {

	searchRes, err := s.client.Search(name, spotify.SearchTypeArtist)
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
			Name: spotifyArtist.Name,
			Genres: spotifyArtist.Genres,
			ArtworkUrl: imageUrl,
			Url: url,
		}

		res = append(res, artist)
	}

	return res, nil
}

func (s *spotifyStreamingService) SearchAlbum(name string) (res []streamingService.Album, err error) {

	searchRes, err := s.client.Search(name, spotify.SearchTypeAlbum)
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
			Name: spotifyAlbum.Name,
			ArtistName: artistName(spotifyAlbum.Artists),
			ArtworkUrl: imageUrl,
			Url: url,
		}

		res = append(res, album)
	}

	return res, nil
}

func (s *spotifyStreamingService) SearchSong(name string) (res []streamingService.Song, err error) {

	searchRes, err := s.client.Search(name, spotify.SearchTypeTrack)
	if err != nil {
		return res, err
	}

	for _, spotifySong := range searchRes.Tracks.Tracks {

		url := spotifySong.ExternalURLs["spotify"]

		song := streamingService.Song{
			Name: spotifySong.Name,
			ArtistName: artistName(spotifySong.Artists),
			AlbumName: spotifySong.Album.Name,
			Url: url,
		}

		res = append(res, song)
	}

	return res, nil
}

func artistName(artists []spotify.SimpleArtist) string {

	var name string
	if len(artists) > 0 {
		for i, artist := range artists {
			if i > 0 && i == len(artists) - 1 {
				name += ", "
			}

			name += artist.Name
		}
	}

	return name
}