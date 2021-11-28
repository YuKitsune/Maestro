package model

import (
	"crypto/sha1"
	"fmt"
	"net/url"
	"strings"
)

type Album struct {
	Name     string
	ArtistName string
	ArtworkLink *url.URL

	Hash ThingHash
	Source StreamingServiceKey
	Market Market
	Link *url.URL
}

func NewAlbum(name string, artistName string, artworkLink *url.URL, source StreamingServiceKey, market Market, link *url.URL) *Album {

	str := fmt.Sprintf("%s_%s", strings.ToLower(artistName), strings.ToLower(name))
	hash := ThingHash(sha1.New().Sum([]byte(str)))

	return &Album{
		Name:       name,
		ArtistName: artistName,
		ArtworkLink: artworkLink,
		Hash: 		hash,
		Source:     source,
		Market:     market,
		Link:       link,
	}
}

func (a *Album) Type() ThingType {
	return AlbumThing
}

func (a *Album) GetHash() ThingHash {
	return a.Hash
}

func (a *Album) GetSource() StreamingServiceKey {
	return a.Source
}

func (a *Album) GetMarket()Market {
	return a.Market
}

func (a *Album) GetLink() *url.URL {
	return a.Link
}
