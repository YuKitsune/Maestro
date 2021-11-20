package streamingService

type Artist struct {
	Name string
	Genres []string
	Url string
}

type Album struct {
	Name string
	ArtistName string
	ArtworkUrl string
	Url string
}

type Song struct {
	Name string
	ArtistName string
	AlbumName string
	Url string
}

type StreamingService interface {
	SearchArtist(name string) ([]*Artist, error)
	SearchAlbum(name string) ([]*Album, error)
	SearchSong(name string) ([]*Song, error)
}
