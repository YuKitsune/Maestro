package deezer

type searchArtistResponse struct {
	Data []Artist
}

type Artist struct {
	Name    string
	Link    string
	Picture string
}

type searchAlbumResponse struct {
	Data []Album
}

type Album struct {
	Title  string
	Link   string
	Cover  string
	Artist Artist
}

type searchTrackResponse struct {
	Data []Track
}

type Track struct {
	Title  string
	Link   string
	Position int `json:"track_position"`
	Artist Artist
	Album  Album
}
