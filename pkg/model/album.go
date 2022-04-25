package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

const AlbumCollectionName = "albums"

type Album struct {
	AlbumId     string
	Name        string
	ArtistNames []string
	ArtworkLink string

	Source StreamingServiceKey
	Market Market
	Link   string
}

func NewAlbum(name string, artistNames []string, artworkLink string, source StreamingServiceKey, market Market, link string) *Album {
	return &Album{
		Name:        name,
		ArtistNames: artistNames,
		ArtworkLink: artworkLink,
		Source:      source,
		Market:      market,
		Link:        link,
	}
}

func (a *Album) GetArtworkLink() string {
	return a.ArtworkLink
}

func (a *Album) GetSource() StreamingServiceKey {
	return a.Source
}

func (a *Album) GetMarket() Market {
	return a.Market
}

func (a *Album) GetLink() string {
	return a.Link
}

func (a *Album) GetLabel() string {
	return fmt.Sprintf("%s (%s)", a.Name, strings.Join(a.ArtistNames, ", "))
}

func UnmarshalAlbum(raw bson.Raw) (*Album, error) {
	var album *Album
	if err := bson.Unmarshal(raw, &album); err != nil {
		return nil, err
	}

	return album, nil
}

func UnmarshalAlbumFromCursor(ctx context.Context, cur *mongo.Cursor) ([]*Album, error) {
	var albums []*Album

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		album, err := UnmarshalAlbum(cur.Current)
		if err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}

	return albums, nil
}
