package db

import (
	"context"
	"github.com/yukitsune/maestro/pkg/model"
)

type Repository interface {
	AddArtist(ctx context.Context, artists []*model.Artist) (int, error)
	GetArtistsById(ctx context.Context, id string) ([]*model.Artist, error)
	GetArtistByLink(ctx context.Context, link string) (*model.Artist, error)

	AddAlbum(ctx context.Context, albums []*model.Album) (int, error)
	GetAlbumsById(ctx context.Context, id string) ([]*model.Album, error)
	GetAlbumByLink(ctx context.Context, link string) (*model.Album, error)

	AddTracks(ctx context.Context, tracks []*model.Track) (int, error)
	GetTracksByLegacyId(ctx context.Context, id string) ([]*model.Track, error)
	GetTracksByIsrc(ctx context.Context, isrc string) ([]*model.Track, error)
	GetTrackByLink(ctx context.Context, link string) (*model.Track, error)

	GetByLink(ctx context.Context, link string) (model.Type, interface{}, error)
}
