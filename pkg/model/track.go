package model

import (
	"fmt"
	"maestro/pkg/hasher"
	"strings"
)

type Track struct {
	Name     string
	ArtistNames []string
	AlbumName  string

	Hash ThingHash
	Source StreamingServiceKey
	ThingType ThingType
	Market Market
	Link string
}

func NewTrack(name string, artistNames []string, albumName string, source StreamingServiceKey, market Market, link string) *Track {

	str := fmt.Sprintf("%s_%s_%s", strings.ToLower(strings.Join(artistNames, "&")), strings.ToLower(albumName), strings.ToLower(name))
	hash := ThingHash(hasher.NewSha1Hasher().ComputeHash(str))

	return &Track{
		name,
		artistNames,
		albumName,
		hash,
		source,
		TrackThing,
		market,
		link,
	}
}

func (t *Track) Type() ThingType {
	return t.ThingType
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

func (t *Track) GetLink() string {
	return t.Link
}
