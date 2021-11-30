package model

import (
	"fmt"
	"maestro/pkg/hasher"
	"sort"
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

	sort.Strings(artistNames)
	str := strings.ToLower(fmt.Sprintf("%s_%s", strings.Join(artistNames, "&"), name))
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
