package model

import (
	"crypto/sha1"
	"net/url"
	"strings"
)

type Artist struct {
	Name string
	ArtworkLink *url.URL

	Hash ThingHash
	Source StreamingServiceKey
	Market Market
	Link *url.URL
}

func NewArtist(name string, artworkLink *url.URL, source StreamingServiceKey, market Market, link *url.URL) *Artist {

	hash := ThingHash(sha1.New().Sum([]byte(strings.ToLower(name))))

	return &Artist{
		Name:   name,
		ArtworkLink: artworkLink,
		Hash: hash,
		Source: source,
		Market: market,
		Link:   link,
	}
}

func (a *Artist) Type() ThingType {
	return ArtistThing
}

func (a *Artist) GetHash() ThingHash {
	return a.Hash
}

func (a *Artist) GetSource() StreamingServiceKey {
	return a.Source
}

func (a *Artist) GetMarket()Market {
	return a.Market
}

func (a *Artist) GetLink() *url.URL {
	return a.Link
}