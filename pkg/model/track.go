package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

const TrackCollectionName = "tracks"

type Track struct {
	Isrc        string
	Name        string
	ArtistNames []string
	AlbumName   string
	ArtworkLink string

	Source StreamingServiceKey
	Market Market
	Link   string
}

func NewTrack(isrc string, name string, artistNames []string, albumName string, artworkLink string, source StreamingServiceKey, market Market, link string) *Track {
	return &Track{
		Isrc:        isrc,
		Name:        name,
		ArtistNames: artistNames,
		AlbumName:   albumName,
		ArtworkLink: artworkLink,
		Source:      source,
		Market:      market,
		Link:        link,
	}
}

func (t *Track) GetArtworkLink() string {
	return t.ArtworkLink
}

func (t *Track) GetSource() StreamingServiceKey {
	return t.Source
}

func (t *Track) GetMarket() Market {
	return t.Market
}

func (t *Track) GetLink() string {
	return t.Link
}

func (t *Track) GetLabel() string {
	return fmt.Sprintf("%s - %s", strings.Join(t.ArtistNames, ", "), t.Name)
}

func UnmarshalTrack(raw bson.Raw) (*Track, error) {
	var track *Track
	if err := bson.Unmarshal(raw, &track); err != nil {
		return nil, err
	}

	return track, nil
}

func UnmarshalTracksFromCursor(ctx context.Context, cur *mongo.Cursor) ([]Track, error) {
	var tracks []Track

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		track, err := UnmarshalTrack(cur.Current)
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, *track)
	}

	return tracks, nil
}
