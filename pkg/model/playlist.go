package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const PlaylistCollectionName = "playlist"

type Playlist struct {
	// PlaylistId is Maestros own ID for the playlist
	PlaylistId string

	// ServicePlaylistId is the ID of the playlist on it's source service
	ServicePlaylistId string

	Name        string
	ArtworkLink string

	Source StreamingServiceKey
	Link   string
}

func NewPlaylist(servicePlaylistId string, name string, artworkLink string, source StreamingServiceKey, link string) *Playlist {
	return &Playlist{
		ServicePlaylistId: servicePlaylistId,
		Name:              name,
		ArtworkLink:       artworkLink,
		Source:            source,
		Link:              link,
	}
}

func (p *Playlist) GetArtworkLink() string {
	return p.ArtworkLink
}

func (p *Playlist) GetSource() StreamingServiceKey {
	return p.Source
}

func (p *Playlist) GetLink() string {
	return p.Link
}

func UnmarshalPlaylist(raw bson.Raw) (*Playlist, error) {
	var playlist *Playlist
	if err := bson.Unmarshal(raw, &playlist); err != nil {
		return nil, err
	}

	return playlist, nil
}

func UnmarshalPlaylistFromCursor(ctx context.Context, cur *mongo.Cursor) ([]*Playlist, error) {
	var playlists []*Playlist

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		playlist, err := UnmarshalPlaylist(cur.Current)
		if err != nil {
			return nil, err
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func PlaylistsToHasStreamingServiceSlice(playlists []*Playlist) []HasStreamingService {
	var s []HasStreamingService
	for _, playlist := range playlists {
		s = append(s, playlist)
	}

	return s
}
