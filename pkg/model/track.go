package model

import (
	"crypto/sha1"
	"fmt"
	"net/url"
	"strings"
)

type Track struct {
	Name     string
	ArtistName string
	AlbumName  string

	Hash ThingHash
	Source StreamingServiceKey
	Market Market
	Link *url.URL
}

func NewTrack(name string, artistName string, albumName string, source StreamingServiceKey, market Market, link *url.URL) *Track {

	str := fmt.Sprintf("%s_%s_%s", strings.ToLower(artistName), strings.ToLower(albumName), strings.ToLower(name))
	hash := ThingHash(sha1.New().Sum([]byte(str)))

	return &Track{
		name,
		artistName,
		albumName,
		hash,
		source,
		market,
		link,
	}
}

func (t *Track) Type() ThingType {
	return TrackThing
}

func (t *Track) GetHash() ThingHash {
	return t.Hash
}

func (t *Track) GetSource() StreamingServiceKey {
	return t.Source
}

func (t *Track) GetMarket()Market {
	return t.Market
}

func (t *Track) GetLink() *url.URL {
	return t.Link
}
