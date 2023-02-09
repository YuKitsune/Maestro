package streamingservice

import (
	"github.com/yukitsune/maestro/pkg/config"
	"github.com/yukitsune/maestro/pkg/model"
)

type StreamingServices map[model.StreamingServiceType]StreamingService

type StreamingService interface {
	Config() config.Service
	LinkBelongsToService(link string) bool
	CleanLink(link string) string

	SearchArtist(artist *model.Artist) (*model.Artist, bool, error)

	SearchAlbum(album *model.Album) (*model.Album, bool, error)

	SearchTrack(song *model.Track) (*model.Track, bool, error)
	GetTrackByIsrc(isrc string) (*model.Track, bool, error)

	GetFromLink(link string) (model.Type, interface{}, error)
}
