package appleMusic

type SearchResponse struct {
	Results SearchResult
}

type ArtistsResult struct {
	Data []Artist
}

type AlbumResult struct {
	Data []Album
}

type SongResult struct {
	Data []Song
}

type Artist struct {
	Id         string
	Attributes ArtistAttributes //The attributes for the artist.
}

type ArtistAttributes struct {
	GenreNames []string //(Required) The names of the genres associated with this artist.
	Name       string   //(Required) The localized name of the artist.
	Url        string   //(Required) The URL for sharing an artist in the iTunes Store.
}

type Song struct {
	Id         string
	Attributes SongAttributes //The attributes for the song.
}

type SongAttributes struct {
	AlbumName  string //(Required) The name of the album the song appears on.
	ArtistName string //(Required) The artist’s name.
	Name       string //(Required) The localized name of the song.
	Url        string //(Required) The URL for sharing a song in the iTunes Store.
}

type Artwork struct {
	BgColor    string
	Height     int
	Width      int
	TextColor1 string
	TextColor2 string
	TextColor3 string
	TextColor4 string
	Url        string
}

type Album struct {
	Href       string
	Id         string
	Attributes AlbumAttributes //The attributes for the album.
}

type AlbumAttributes struct {
	AlbumName  string  //(Required) The name of the album the music video appears on.
	ArtistName string  //(Required) The artist’s name.
	Artwork    Artwork //The album artwork.
	Name       string  //(Required) The localized name of the album.
	Url        string
}

type QueryParams struct {
	Term  string
	Types []string
}

type SearchResult struct {
	Artists *ArtistsResult
	Albums  *AlbumResult
	Songs   *SongResult
}
