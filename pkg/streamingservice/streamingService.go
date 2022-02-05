package streamingservice

import (
	"fmt"
	"github.com/yukitsune/maestro/pkg/model"
)

type StreamingService interface {
	Key() model.StreamingServiceKey
	LinkBelongsToService(link string) bool
	SearchArtist(artist *model.Artist) (*model.Artist, error)
	SearchAlbum(album *model.Album) (*model.Album, error)
	SearchSong(song *model.Track) (*model.Track, error)
	SearchFromLink(link string) (model.Thing, error)
	CleanLink(link string) string
}

func SearchThing(ss StreamingService, thing model.Thing) (model.Thing, error) {

	// Todo: Name normalisation
	// 	Some services name things differently, making searches much more difficult.
	//	Example:
	//		https://open.spotify.com/album/0HVuB4xTfASjULQAfRRS3s
	//		Spotify:
	//			Album Name: Eat Sleep Dance
	//			Artists: ["電音部", "Moe Shop"]
	//		Apple Music:
	//			Album Name: Eat Sleep Dance (feat. Moe Shop)
	//			Artists: ["DENONBU", "Inubousaki Shian (CV: Rena Hasegawa)"]

	// Different languages for artist names, and different naming formats.
	// Apple Music links will result in a more strict search.
	// Should normalise these names to our own format, so we can more easily search

	switch t := thing.(type) {
	case *model.Artist:
		return ss.SearchArtist(t)

	case *model.Album:
		return ss.SearchAlbum(t)

	case *model.Track:
		return ss.SearchSong(t)

	default:
		return nil, fmt.Errorf("unknown type %T", thing)
	}
}

func ForEachStreamingService(services []StreamingService, fn func(StreamingService) error) error {
	for _, service := range services {
		err := fn(service)
		if err != nil {
			return err
		}
	}

	return nil
}
