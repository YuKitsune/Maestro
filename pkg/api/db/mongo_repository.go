package db

import (
	"context"
	"github.com/yukitsune/maestro/pkg/metrics"
	"github.com/yukitsune/maestro/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	db  *mongo.Database
	rec metrics.Recorder
}

func NewMongoRepository(db *mongo.Database, rec metrics.Recorder) Repository {
	return &mongoRepository{db, rec}
}

func (m *mongoRepository) AddArtist(ctx context.Context, artists []model.Artist) (int, error) {
	go m.rec.CountDatabaseCall()

	coll := m.db.Collection(model.ArtistCollectionName)
	insertRes, err := coll.InsertMany(ctx, artistsToInterfaces(artists))
	if err != nil {
		return 0, err
	}

	return len(insertRes.InsertedIDs), nil
}

func (m *mongoRepository) GetArtistsById(ctx context.Context, id string) ([]model.Artist, error) {
	go m.rec.CountDatabaseCall()

	coll := m.db.Collection(model.ArtistCollectionName)
	cur, err := coll.Find(ctx, bson.D{
		{"id", id}, // Todo: Mongo ID?
	})

	artists, err := model.UnmarshalArtistFromCursor(ctx, cur)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (m *mongoRepository) GetArtistByLink(ctx context.Context, link string) (*model.Artist, error) {
	go m.rec.CountDatabaseCall()

	// Find an artist with a matching link
	var foundArtist *model.Artist
	coll := m.db.Collection(model.ArtistCollectionName)
	res := coll.FindOne(ctx, bson.D{{"link", link}})
	err := res.Err()

	// No matches? Error time
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	} else {
		// Found matches? Gimmi gimmi!
		raw, err := res.DecodeBytes()
		if err != nil {
			return nil, err
		}

		foundArtist, err = model.UnmarshalArtist(raw)
		if err != nil {
			return nil, err
		}
	}

	return foundArtist, nil
}

func (m *mongoRepository) UpdateArtists(ctx context.Context, artists []model.Artist) (int, error) {
	go m.rec.CountDatabaseCall()
	panic("not implemented!")
}

func (m *mongoRepository) AddAlbum(ctx context.Context, albums []model.Album) (int, error) {
	go m.rec.CountDatabaseCall()

	coll := m.db.Collection(model.AlbumCollectionName)
	insertRes, err := coll.InsertMany(ctx, albumsToInterfaces(albums))
	if err != nil {
		return 0, err
	}

	return len(insertRes.InsertedIDs), nil
}

func (m *mongoRepository) GetAlbumsById(ctx context.Context, id string) ([]model.Album, error) {
	go m.rec.CountDatabaseCall()

	coll := m.db.Collection(model.AlbumCollectionName)
	cur, err := coll.Find(ctx, bson.D{
		{"id", id}, // Todo: Mongo ID?
	})

	albums, err := model.UnmarshalAlbumFromCursor(ctx, cur)
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (m *mongoRepository) GetAlbumByLink(ctx context.Context, link string) (*model.Album, error) {
	go m.rec.CountDatabaseCall()

	// Find an album with a matching link
	var foundAlbum *model.Album
	coll := m.db.Collection(model.AlbumCollectionName)
	res := coll.FindOne(ctx, bson.D{{"link", link}})
	err := res.Err()

	// No matches? Error time
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	} else {
		// Found matches? Gimmi gimmi!
		raw, err := res.DecodeBytes()
		if err != nil {
			return nil, err
		}

		foundAlbum, err = model.UnmarshalAlbum(raw)
		if err != nil {
			return nil, err
		}
	}

	return foundAlbum, nil
}

func (m *mongoRepository) UpdateAlbums(ctx context.Context, albums []model.Album) (int, error) {
	go m.rec.CountDatabaseCall()
	panic("not implemented!")
}

func (m *mongoRepository) AddTracks(ctx context.Context, tracks []model.Track) (int, error) {
	go m.rec.CountDatabaseCall()

	coll := m.db.Collection(model.TrackCollectionName)
	insertRes, err := coll.InsertMany(ctx, tracksToInterfaces(tracks))
	if err != nil {
		return 0, err
	}

	return len(insertRes.InsertedIDs), nil
}

func (m *mongoRepository) GetTracksByIsrc(ctx context.Context, isrc string) ([]model.Track, error) {
	go m.rec.CountDatabaseCall()

	coll := m.db.Collection(model.TrackCollectionName)
	cur, err := coll.Find(ctx, bson.D{
		{"isrc", isrc},
	})

	tracks, err := model.UnmarshalTracksFromCursor(ctx, cur)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (m *mongoRepository) GetTrackByLink(ctx context.Context, link string) (*model.Track, error) {
	go m.rec.CountDatabaseCall()

	// Find a track with a matching link
	var foundTrack *model.Track
	coll := m.db.Collection(model.TrackCollectionName)
	res := coll.FindOne(ctx, bson.D{{"link", link}})
	err := res.Err()

	// No matches? Error time
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	} else {
		// Found matches? Gimmi gimmi!
		raw, err := res.DecodeBytes()
		if err != nil {
			return nil, err
		}

		foundTrack, err = model.UnmarshalTrack(raw)
		if err != nil {
			return nil, err
		}
	}

	return foundTrack, nil
}

func (m *mongoRepository) UpdateTracks(ctx context.Context, tracks []model.Track) (int, error) {
	go m.rec.CountDatabaseCall()
	panic("not implemented!")
}

func artistsToInterfaces(artists []model.Artist) []interface{} {
	var s []interface{}
	for _, artist := range artists {
		s = append(s, artist)
	}

	return s
}

func albumsToInterfaces(albums []model.Album) []interface{} {
	var s []interface{}
	for _, album := range albums {
		s = append(s, album)
	}

	return s
}

func tracksToInterfaces(tracks []model.Track) []interface{} {
	var s []interface{}
	for _, track := range tracks {
		s = append(s, track)
	}

	return s
}
