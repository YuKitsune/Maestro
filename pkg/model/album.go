package model

import (
	"fmt"
	"maestro/pkg/hasher"
	"strings"
)

type Album struct {
	Name     string
	ArtistNames []string
	ArtworkLink string

	Hash ThingHash
	Source StreamingServiceKey
	ThingType ThingType
	Market Market
	Link string
}

func NewAlbum(name string, artistNames []string, artworkLink string, source StreamingServiceKey, market Market, link string) *Album {

	str := fmt.Sprintf("%s_%s", strings.ToLower(strings.Join(artistNames, "&")), strings.ToLower(name))
	hash := ThingHash(hasher.NewSha1Hasher().ComputeHash(str))

	return &Album{
		name,
		artistNames,
		artworkLink,
		hash,
		source,
		AlbumThing,
		market,
		link,
	}
}

func (a *Album) Type() ThingType {
	return a.ThingType
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

func (a *Album) GetLink() string {
	return a.Link
}
