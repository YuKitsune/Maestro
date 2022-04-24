package db

import (
	"context"
	"github.com/yukitsune/maestro/pkg/model"
)

type Repository interface {
	AddArtist(ctx context.Context, artists []model.Artist) (int, error)
	GetArtistsById(ctx context.Context, id string) ([]model.Artist, error)
	GetArtistByLink(ctx context.Context, link string) (*model.Artist, error)
	UpdateArtists(ctx context.Context, artists []model.Artist) (int, error)

	AddAlbum(ctx context.Context, albums []model.Album) (int, error)
	GetAlbumsById(ctx context.Context, id string) ([]model.Album, error)
	GetAlbumByLink(ctx context.Context, link string) (*model.Album, error)
	UpdateAlbums(ctx context.Context, albums []model.Album) (int, error)

	AddTracks(ctx context.Context, tracks []model.Track) (int, error)
	GetTracksByIsrc(ctx context.Context, isrc string) ([]model.Track, error)
	GetTrackByLink(ctx context.Context, link string) (*model.Track, error)
	UpdateTracks(ctx context.Context, tracks []model.Track) (int, error)
}
