package main

import (
	"fmt"
	"maestro/pkg/streamingService/appleMusic"
	"maestro/pkg/streamingService/spotify"
)

func main() {

	am := appleMusic.NewAppleMusicStreamingService(todo)
	sp := spotify.NewSpotifyStreamingService(todo)

	fmt.Println("artist, album, or song?")

	var typ string
	fmt.Scanln(&typ)

	switch typ {
	case "artist":
		var term string
		getTerm(&term)

		// Todo: Apple Music giving empty results for some reason...
		amArtist, err := am.SearchArtist(term)
		if err != nil {
			panic(err)
		}

		spArtist, err := sp.SearchArtist(term)
		if err != nil {
			panic(err)
		}

		fmt.Println("Apple Music:")
		for _, artist := range amArtist {
			fmt.Printf("\t%s\t%s\n", artist.Name, artist.Url)
		}

		fmt.Println("Spotify:")
		for _, artist := range spArtist {
			fmt.Printf("\t%s\t%s\n", artist.Name, artist.Url)
		}

		break

	case "album":
		break

	case "song":

		break

	default:
		fmt.Println("Huh???")
		break
	}

}

func getTerm(v interface{}) {
	fmt.Printf("search: ")
	fmt.Scanln(v)
}