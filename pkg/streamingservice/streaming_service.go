package streamingservice

import (
	"github.com/yukitsune/maestro/pkg/model"
)

type StreamingServices map[model.StreamingServiceKey]StreamingService

type StreamingService interface {
	LinkBelongsToService(link string) bool
	CleanLink(link string) string

	SearchArtist(artist *model.Artist) (*model.Artist, bool, error)

	SearchAlbum(album *model.Album) (*model.Album, bool, error)

	SearchTrack(song *model.Track) (*model.Track, bool, error)
	GetTrackByIsrc(isrc string) (*model.Track, bool, error)

	GetPlaylistById(id string) (*model.Playlist, bool, error)
	GetPlaylistTracksById(id string) ([]*model.Track, bool, error)

	GetFromLink(link string) (model.Type, interface{}, error)
}
