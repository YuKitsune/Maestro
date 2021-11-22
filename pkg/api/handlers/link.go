package handlers

import "net/http"

func HandleLink(w http.ResponseWriter, r *http.Request) {

	requestedUrl := "https://music.apple.com/au/album/i-gotta/1453467630?i=1453467634"

	// 1. Check links table for a matching record

	// If a record exists:
	//	Compile records
	// 	If we're aware of a streaming service, but there are no results for that service (E.g: A new service added after the result was stored)
	//		1. Get details from any remaining streaming service APIs
	//  	2. Store the new details in the database
	//		2. Add the new details to our result set

	// If no records exist:
	// 	1. Get details from the streaming service API
	// 	2. Search other streaming services for similar details
	// 	3. Store the best matches in the database

	//  Return the results for each streaming service
}

func HandleFlagLink(w http.ResponseWriter, r *http.Request) {

	// Todo: Pick one???
	// Option 1: Delete the records
	// Option 2: Delete the records after N number of flags

}
