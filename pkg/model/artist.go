package model

import (
	"maestro/pkg/hasher"
)

type Artist struct {
	Name string
	ArtworkLink string

	Hash ThingHash
	Source StreamingServiceKey
	ThingType ThingType
	Market Market
	Link string
}

func NewArtist(name string, artworkLink string, source StreamingServiceKey, market Market, link string) *Artist {

	hash := ThingHash(hasher.NewSha1Hasher().ComputeHash(name))

	return &Artist{
		name,
		artworkLink,
		hash,
		source,
		ArtistThing,
		market,
		link,
	}
}

func (a *Artist) Type() ThingType {
	return a.ThingType
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

func (a *Artist) GetLink() string {
	return a.Link
}