package streamingService

type Region string

func RegionToString(r Region) string {
	var rv interface{}
	rv = r
	rs := rv.(string)
	return rs
}

type Thing interface {
	GetName() string
	GetUrl() string
}

type Artist struct {
	Name       string
	ArtworkUrl string
	Url        string
}

func (a *Artist) GetName() string {
	return a.Name
}

func (a *Artist) GetUrl() string {
	return a.Url
}

type Album struct {
	Name       string
	ArtistName string
	ArtworkUrl string
	Url        string
}

func (a *Album) GetName() string {
	return a.Name
}

func (a *Album) GetUrl() string {
	return a.Url
}

type Song struct {
	Name       string
	ArtistName string
	AlbumName  string
	Url        string
}

func (a *Song) GetName() string {
	return a.Name
}

func (a *Song) GetUrl() string {
	return a.Url
}

type StreamingService interface {
	Name() string
	SearchArtist(name string, region Region) ([]Artist, error)
	SearchAlbum(name string, region Region) ([]Album, error)
	SearchSong(name string, region Region) ([]Song, error)
	SearchFromLink(link string) (Thing, error)
}
