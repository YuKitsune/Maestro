package handlers

import (
	"context"
	"fmt"
	"github.com/yukitsune/maestro/pkg/db"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/maestro/pkg/log"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingservice"
)

func GetLinkHandler(serviceProvider streamingservice.ServiceProvider, repo db.Repository, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLogger, err := log.ForRequest(logger, r)
		if err != nil {
			Error(w, err)
			return
		}

		vars := mux.Vars(r)
		reqLink, ok := vars["link"]
		if !ok {
			BadRequest(w, "missing parameter \"link\"")
			return
		}

		u, err := url.Parse(reqLink)
		if err != nil || u == nil {
			BadRequestf(w, "couldn't parse the given link: %s", reqLink)
			return
		}

		if !u.IsAbs() {
			BadRequest(w, "given link must be absolute")
			return
		}

		res, err := findForLink(r.Context(), reqLink, serviceProvider, repo, reqLogger)
		if err != nil {
			Error(w, err)
			return
		}

		if res == nil || !res.HasResults() {
			NotFound(w, "could not find anything")
			return
		}

		Response(w, res, http.StatusOK)
	}
}

func findForLink(ctx context.Context, link string, serviceProvider streamingservice.ServiceProvider, repo db.Repository, logger *logrus.Entry) (*Result, error) {
	svcs, err := serviceProvider.ListServices()
	if err != nil {
		return nil, err
	}

	// Trim service-specific stuff from the link
	for _, svc := range svcs {
		link = svc.CleanLink(link)
	}

	logger = logger.WithField("link", link)

	// Search the database for an existing thing with the given link
	typ, dbRes, err := repo.GetByLink(ctx, link)
	if err != nil {
		return nil, err
	}

	switch typ {
	case model.ArtistType:
		artist := dbRes.(*model.Artist)
		res, err := findForExistingArtist(ctx, artist, svcs, repo, logger)
		return res, err

	case model.AlbumType:
		album := dbRes.(*model.Album)
		res, err := findForExistingAlbum(ctx, album, svcs, repo, logger)
		return res, err

	case model.TrackType:
		track := dbRes.(*model.Track)
		res, err := findForExistingTrack(ctx, track, svcs, repo, logger)
		return res, err

	case model.UnknownType:
		res, err := findNewThing(ctx, link, svcs, repo, logger)
		return res, err

	default:
		return nil, fmt.Errorf("unknown type %s", typ)
	}
}

func findForExistingArtist(ctx context.Context, foundArtist *model.Artist, svcs streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (*Result, error) {

	logger = logger.WithField("artist_id", foundArtist.ArtistId)
	logger.Debugln("found an artist")

	res := NewResult(model.ArtistType)

	// Find any related artists based on our artist ID
	existingArtists, err := repo.GetArtistsById(ctx, foundArtist.ArtistId)
	if err != nil {
		return nil, err
	}

	res.AddAll(model.ArtistsToHasStreamingServiceSlice(existingArtists))

	// If we have results for all known services, then we're good to go
	if len(res.Items) == len(svcs) {
		return res, nil
	}

	logger.Debugf("looks like we have some new services since we found this artist (found %d, looking for %d)\n", len(existingArtists), len(svcs))

	// Query the remaining streaming service
	var newArtists []*model.Artist
	for key, svc := range svcs {
		if res.HasResultFor(key) {
			continue
		}

		logger.Debugf("searching %s for artist\n", key)
		artist, found, err := svc.SearchArtist(foundArtist)
		if err != nil {
			logger.Errorf("%s: %s", key, err.Error())
			continue
		}

		if !found {
			logger.Debugf("couldn't find anything for %s", key)
			continue
		}

		artist.ArtistId = foundArtist.ArtistId
		newArtists = append(newArtists, artist)
	}

	// Add the new artists to the database
	if len(newArtists) != 0 {

		n, err := repo.AddArtist(ctx, newArtists)
		if err != nil {
			return nil, err
		}

		logger.Infof("%d new artists added\n", n)
	}

	res.AddAll(model.ArtistsToHasStreamingServiceSlice(newArtists))
	return res, nil
}

func findForExistingAlbum(ctx context.Context, foundAlbum *model.Album, svcs streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (*Result, error) {

	logger = logger.WithField("album_id", foundAlbum.AlbumId)
	logger.Debugln("found an album")

	res := NewResult(model.AlbumType)

	// Find any related albums based on our album ID
	existingAlbums, err := repo.GetAlbumsById(ctx, foundAlbum.AlbumId)
	if err != nil {
		return nil, err
	}

	res.AddAll(model.AlbumToHasStreamingServiceSlice(existingAlbums))

	// If we have results for all known services, then we're good to go
	if len(res.Items) == len(svcs) {
		return res, nil
	}

	logger.Debugf("looks like we have some new services since we found this album (found %d, looking for %d)\n", len(existingAlbums), len(svcs))

	// Query the remaining streaming service
	var newAlbums []*model.Album
	for key, svc := range svcs {
		if res.HasResultFor(key) {
			continue
		}

		logger.Debugf("searching %s for album\n", key)
		album, found, err := svc.SearchAlbum(foundAlbum)
		if err != nil {
			logger.Errorf("%s: %s", key, err.Error())
			continue
		}

		if !found {
			logger.Debugf("couldn't find anything for %s", key)
			continue
		}

		album.AlbumId = foundAlbum.AlbumId
		newAlbums = append(newAlbums, album)
	}

	// Add the new albums to the database
	if len(newAlbums) != 0 {

		n, err := repo.AddAlbum(ctx, newAlbums)
		if err != nil {
			return nil, err
		}

		logger.Infof("%d new albums added\n", n)
	}

	res.AddAll(model.AlbumToHasStreamingServiceSlice(newAlbums))
	return res, nil
}

func findForExistingTrack(ctx context.Context, foundTrack *model.Track, svcs streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (*Result, error) {

	logger = logger.WithField("isrc", foundTrack.Isrc)
	logger.Debugln("found a track")

	res := NewResult(model.TrackType)

	// Find any related track based on the ISRC
	existingTracks, err := repo.GetTracksByIsrc(ctx, foundTrack.Isrc)
	if err != nil {
		return nil, err
	}

	res.AddAll(model.TrackToHasStreamingServiceSlice(existingTracks))

	// If we have results for all known services, then we're good to go
	if len(res.Items) == len(svcs) {
		return res, nil
	}

	logger.Debugf("looks like we have some new services since we found this track (found %d, looking for %d)\n", len(existingTracks), len(svcs))

	// Query the remaining streaming service
	var newTracks []*model.Track
	for key, svc := range svcs {
		if res.HasResultFor(key) {
			continue
		}

		logger.Debugf("searching %s for track\n", key)
		track, found, err := svc.GetTrackByIsrc(foundTrack.Isrc)
		if err != nil {
			logger.Errorf("%s: %s", key, err.Error())
			continue
		}

		if !found {
			logger.Debugf("couldn't find anything for %s", key)
			continue
		}

		newTracks = append(newTracks, track)
	}

	// Add the new tracks to the database
	if len(newTracks) != 0 {

		n, err := repo.AddTracks(ctx, newTracks)
		if err != nil {
			return nil, err
		}

		logger.Infof("%d new tracks added\n", n)
	}

	res.AddAll(model.TrackToHasStreamingServiceSlice(newTracks))
	return res, nil
}

func handleNewArtist(ctx context.Context, newArtist *model.Artist, svcs streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (*Result, error) {

	res := NewResult(model.ArtistType)
	res.Add(newArtist)

	// Create our own ID for the artist
	id := uuid.New().String()
	newArtist.ArtistId = id

	logger = logger.WithField("artist_id", id)
	logger.Debugln("using new artist id")

	newArtists := []*model.Artist{
		newArtist,
	}

	// Query the other streaming services using what we found from the target streaming service
	for key, svc := range svcs {
		if res.HasResultFor(key) {
			continue
		}

		logger.Debugf("searching %s for artist with name %s\n", key, newArtist.Name)
		foundArtist, found, err := svc.SearchArtist(newArtist)
		if err != nil {
			logger.Errorf("%s: %s", key, err.Error())
			continue
		}

		if !found {
			logger.Debugf("couldn't find anything for %s", key)
			continue
		}

		foundArtist.ArtistId = id
		res.Add(foundArtist)
		newArtists = append(newArtists, foundArtist)
	}

	n, err := repo.AddArtist(ctx, newArtists)
	if err != nil {
		return nil, err
	}

	logger.Infof("%d new artists added", n)

	return res, nil
}

func handleNewAlbum(ctx context.Context, newAlbum *model.Album, svcs streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (*Result, error) {

	res := NewResult(model.AlbumType)
	res.Add(newAlbum)

	// Create our own ID for the album
	id := uuid.New().String()
	newAlbum.AlbumId = id

	logger = logger.WithField("album_id", id)
	logger.Debugln("using new album id")

	newAlbums := []*model.Album{
		newAlbum,
	}

	// Query the other streaming services using what we found from the target streaming service
	for key, svc := range svcs {
		if res.HasResultFor(key) {
			continue
		}

		logger.Debugf("searching %s for album with name %s\n", key, newAlbum.Name)
		foundAlbum, found, err := svc.SearchAlbum(newAlbum)
		if err != nil {
			logger.Errorf("%s: %s", key, err.Error())
			continue
		}

		if !found {
			logger.Debugf("couldn't find anything for %s", key)
			continue
		}

		foundAlbum.AlbumId = id
		res.Add(foundAlbum)
		newAlbums = append(newAlbums, foundAlbum)
	}

	n, err := repo.AddAlbum(ctx, newAlbums)
	if err != nil {
		return nil, err
	}

	logger.Infof("%d new albums added", n)

	return res, nil
}

func handleNewTrack(ctx context.Context, newTrack *model.Track, svcs streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (*Result, error) {

	res := NewResult(model.TrackType)
	res.Add(newTrack)

	logger = logger.WithField("isrc", newTrack.Isrc)

	newTracks := []*model.Track{
		newTrack,
	}

	// Query the other streaming services using what we found from the target streaming service
	for key, svc := range svcs {
		logger.Debugf("searching %s for track with name %s\n", key, newTrack.Name)
		foundTrack, found, err := svc.GetTrackByIsrc(newTrack.Isrc)
		if err != nil {
			logger.Errorf("%s: %s", svc, err.Error())
			continue
		}

		if !found {
			logger.Debugf("couldn't find anything for %s", key)
			continue
		}

		res.Add(foundTrack)
		newTracks = append(newTracks, foundTrack)
	}

	n, err := repo.AddTracks(ctx, newTracks)
	if err != nil {
		return nil, err
	}

	logger.Infof("%d new tracks added", n)

	return res, nil
}

func findNewThing(ctx context.Context, link string, svcs streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (*Result, error) {

	logger.Debugln("looks like this is a new thing")

	// No links found, query the streaming service and find the same entry on other services
	var targetKey model.StreamingServiceType
	var targetService streamingservice.StreamingService
	otherServices := make(streamingservice.StreamingServices)

	// Figure out which streaming service the link belongs to
	for key, svc := range svcs {
		if svc.LinkBelongsToService(link) {
			targetKey = key
			targetService = svc
		} else {
			otherServices[key] = svc
		}
	}

	if targetService == nil {
		return nil, fmt.Errorf("couldn't find a streaming service for the given link: %s", link)
	}

	// Query the target streaming service
	logger.Debugf("searching %s\n", targetKey)
	typ, res, err := targetService.GetFromLink(link)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", targetKey, err.Error())
	}

	switch typ {
	case model.ArtistType:
		artist := res.(*model.Artist)
		return handleNewArtist(ctx, artist, otherServices, repo, logger)

	case model.AlbumType:
		album := res.(*model.Album)
		return handleNewAlbum(ctx, album, otherServices, repo, logger)

	case model.TrackType:
		track := res.(*model.Track)
		return handleNewTrack(ctx, track, otherServices, repo, logger)

	case model.UnknownType:
		return nil, fmt.Errorf("could not find anything from %s", targetKey)

	default:
		return nil, fmt.Errorf("unknown type %s", typ)
	}
}
