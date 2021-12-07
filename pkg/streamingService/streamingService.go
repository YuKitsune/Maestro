package streamingService

import (
	"fmt"
	"maestro/pkg/model"
)

// Todo: Not a fan of these models being separate from the ones in the models package
// 	but the shape is inherently different...
// 	The ones in the models package are more of a metadata with many links thing, where as these ones are the links, but
// 	also need said metadata.
// 	Might be better off flattening the model and using some kind of correlation ID / hash instead to relate entities

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
