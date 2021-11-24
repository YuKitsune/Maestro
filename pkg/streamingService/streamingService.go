package streamingService

type Artist struct {
	Name       string
	Genres     []string
	Url        string
	ArtworkUrl string
}

type Album struct {
	Name       string
	ArtistName string
	ArtworkUrl string
	Url        string
}

type Song struct {
	Name       string
	ArtistName string
	AlbumName  string
	Url        string
}

type StreamingService interface {
	Name() string
	SearchArtist(name string) ([]Artist, error)
	SearchAlbum(name string) ([]Album, error)
	SearchSong(name string) ([]Song, error)
}
