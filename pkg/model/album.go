package model

import (
	"fmt"
	"strings"
)

type Album struct {
	Name        string
	ArtistNames []string
	ArtworkLink string

	GroupID   ThingGroupID
	Source    StreamingServiceKey
	ThingType ThingType
	Market    Market
	Link      string
}

func NewAlbum(name string, artistNames []string, artworkLink string, source StreamingServiceKey, market Market, link string) *Album {
	return &Album{
		Name:        name,
		ArtistNames: artistNames,
		ArtworkLink: artworkLink,
		Source:      source,
		ThingType:   AlbumThing,
		Market:      market,
		Link:        link,
	}
}

func (a *Album) Type() ThingType {
	return a.ThingType
}

func (a *Album) GetArtworkLink() string {
	return a.ArtworkLink
}

func (a *Album) GetGroupID() ThingGroupID {
	return a.GroupID
}

func (a *Album) SetGroupID(groupID ThingGroupID) {
	a.GroupID = groupID
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
