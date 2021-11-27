package streamingService

import (
	"fmt"
	"maestro/pkg/model"
	"net/url"
)

// Todo: Not a fan of these models being separate from the ones in the models package
// 	but the shape is inherently different...
// 	The ones in the models package are more of a metadata with many links thing, where as these ones are the links, but
// 	also need said metadata.
// 	Might be better off flattening the model and using some kind of correlation ID / hash instead to relate entities

type Thing interface {
	GetName() string
	GetUrl() string
	GetMarket() model.Market
}

type Artist struct {
	Name       string
	ArtworkUrl string
	Market     model.Market
	Url        string
}

func (a *Artist) GetName() string {
	return a.Name
}

func (a *Artist) GetUrl() string {
	return a.Url
}

func (a *Artist) GetMarket() model.Market {
	return a.Market
}

type Album struct {
	Name       string
	ArtistName string
	ArtworkUrl string
	Market     model.Market
	Url        string
}

func (a *Album) GetName() string {
	return a.Name
}

func (a *Album) GetUrl() string {
	return a.Url
}

func (a *Album) GetMarket() model.Market {
	return a.Market
}

type Song struct {
	Name       string
	ArtistName string
	AlbumName  string
	Number 	   int
	Market     model.Market
	Url        string
}

func (a *Song) GetName() string {
	return a.Name
}

func (a *Song) GetUrl() string {
	return a.Url
}

func (a *Song) GetMarket() model.Market {
	return a.Market
}

type StreamingService interface {
	Name() model.StreamingServiceKey
	LinkBelongsToService(link string) bool
	SearchArtist(artist *Artist) (*Artist, error)
	SearchAlbum(album *Album) (*Album, error)
	SearchSong(song *Song) (*Song, error)
	SearchFromLink(link string) (Thing, error)
}

func SearchThing(ss StreamingService, thing Thing) (Thing, error) {
	switch t := thing.(type) {
	case *Artist:
		return ss.SearchArtist(t)

	case *Album:
		return ss.SearchAlbum(t)

	case *Song:
		return ss.SearchSong(t)

	default:
		return nil, fmt.Errorf("unknown type %T", thing)
	}
}

func ForEachStreamingService(services []StreamingService, fn func(StreamingService) error) error {
	for _, service := range services {
		err := fn(service)
		if err != nil {
			return err
		}
	}

	return nil
}

func ConvertToModel(service StreamingService, thing Thing) (model.Thing, error) {

	links := make(model.Links)
	url, err := url.Parse(thing.GetUrl())
	if err != nil {
		return nil, err
	}

	links[service.Name()] = model.Link {
		Market: thing.GetMarket(),
		Url:    url,
	}

	switch t := thing.(type) {
	case *Artist:
		return &model.Artist{
			Name:  t.Name,
			Links: links,
		}, nil

	case *Album:
		return &model.Album{
			Name:     t.Name,
			ArtistId: "", // Todo
			Links:    links,
		}, nil

	case *Song:
		return &model.Track{
			Name:     t.Name,
			ArtistId: "", // Todo
			AlbumId:  "", // Todo
			Number:   t.Number,
			Links:    links,
		}, nil

	default:
		return nil, fmt.Errorf("unknown type %T", thing)
	}
}