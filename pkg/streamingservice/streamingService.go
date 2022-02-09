package streamingservice

import (
	"fmt"
	"github.com/yukitsune/maestro/pkg/model"
)

type StreamingService interface {
	Key() model.StreamingServiceKey
	LinkBelongsToService(link string) bool
	SearchArtist(artist *model.Artist) (*model.Artist, bool, error)
	SearchAlbum(album *model.Album) (*model.Album, bool, error)
	SearchSong(song *model.Track) (*model.Track, bool, error)
	SearchFromLink(link string) (model.Thing, bool, error)
	CleanLink(link string) string
}

func SearchThing(ss StreamingService, thing model.Thing) (model.Thing, bool, error) {

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
		return nil, false, fmt.Errorf("unknown type %T", thing)
	}
}
