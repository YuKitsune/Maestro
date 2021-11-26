package streamingService

import "fmt"

type Region string

func RegionToString(r Region) string {
	return fmt.Sprintf("%s", r)
}

type Thing interface {
	GetName() string
	GetUrl() string
	GetRegion() Region
}

type Artist struct {
	Name       string
	ArtworkUrl string
	Region Region
	Url        string
}

func (a *Artist) GetName() string {
	return a.Name
}

func (a *Artist) GetUrl() string {
	return a.Url
}

func (a *Artist) GetRegion() Region {
	return a.Region
}

type Album struct {
	Name       string
	ArtistName string
	ArtworkUrl string
	Region Region
	Url        string
}

func (a *Album) GetName() string {
	return a.Name
}

func (a *Album) GetUrl() string {
	return a.Url
}

func (a *Album) GetRegion() Region {
	return a.Region
}

type Song struct {
	Name       string
	ArtistName string
	AlbumName  string
	Region Region
	Url        string
}

func (a *Song) GetName() string {
	return a.Name
}

func (a *Song) GetUrl() string {
	return a.Url
}

func (a *Song) GetRegion() Region {
	return a.Region
}

type StreamingService interface {
	Name() string
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