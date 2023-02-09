package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const ArtistCollectionName = "artists"

type Artist struct {
	ArtistId    string
	Name        string
	ArtworkLink string

	Source StreamingServiceType
	Market Market
	Link   string
}

func NewArtist(name string, artworkLink string, source StreamingServiceType, market Market, link string) *Artist {
	return &Artist{
		Name:        name,
		ArtworkLink: artworkLink,
		Source:      source,
		Market:      market,
		Link:        link,
	}
}

func (a *Artist) GetArtworkLink() string {
	return a.ArtworkLink
}

func (a *Artist) GetSource() StreamingServiceType {
	return a.Source
}

func (a *Artist) GetMarket() Market {
	return a.Market
}

func (a *Artist) GetLink() string {
	return a.Link
}

func (a *Artist) GetLabel() string {
	return a.Name
}

func UnmarshalArtist(raw bson.Raw) (*Artist, error) {
	var artist *Artist
	if err := bson.Unmarshal(raw, &artist); err != nil {
		return nil, err
	}

	return artist, nil
}

func UnmarshalArtistFromCursor(ctx context.Context, cur *mongo.Cursor) ([]*Artist, error) {
	var artists []*Artist

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		artist, err := UnmarshalArtist(cur.Current)
		if err != nil {
			return nil, err
		}

		artists = append(artists, artist)
	}

	return artists, nil
}

func ArtistsToHasStreamingServiceSlice(artists []*Artist) []HasStreamingService {
	var s []HasStreamingService
	for _, artist := range artists {
		s = append(s, artist)
	}

	return s
}
