package streamingservice

import (
	"github.com/yukitsune/maestro/pkg/model"
)

type StreamingServices []StreamingService

type StreamingService interface {
	Key() model.StreamingServiceKey
	LinkBelongsToService(link string) bool
	CleanLink(link string) string

	SearchArtist(artist *model.Artist) (*model.Artist, bool, error)

	SearchAlbum(album *model.Album) (*model.Album, bool, error)

	SearchTrack(song *model.Track) (*model.Track, bool, error)
	GetTrackByIsrc(isrc string) (*model.Track, bool, error)

	GetFromLink(link string) (model.Type, interface{}, error)
}
