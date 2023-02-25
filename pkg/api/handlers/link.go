package handlers

import (
	"context"
	"fmt"
	"github.com/yukitsune/maestro/pkg/api/responses"
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
			responses.Error(w, err)
			return
		}

		vars := mux.Vars(r)
		reqLink, ok := vars["link"]
		if !ok {
			responses.BadRequest(w, "missing parameter \"link\"")
			return
		}

		u, err := url.Parse(reqLink)
		if err != nil || u == nil {
			responses.BadRequestf(w, "couldn't parse the given link: %s", reqLink)
			return
		}

		if !u.IsAbs() {
			responses.BadRequest(w, "given link must be absolute")
			return
		}

		res, found, err := findForLink(r.Context(), reqLink, serviceProvider, repo, reqLogger)
		if err != nil {
			responses.Error(w, err)
			return
		}

		if !found {
			responses.NotFound(w, "could not find anything")
		}

		responses.Response(w, res, http.StatusOK)
	}
}

func findForLink(ctx context.Context, link string, serviceProvider streamingservice.ServiceProvider, repo db.Repository, logger *logrus.Entry) (any, bool, error) {
	services, err := serviceProvider.ListServices()
	if err != nil {
		return nil, false, err
	}

	// Trim service-specific stuff from the link
	for _, service := range services {
		link = service.CleanLink(link)
	}

	logger = logger.WithField("link", link)

	// Search the database for an existing thing with the given link
	typ, dbRes, err := repo.GetByLink(ctx, link)
	if err != nil {
		return nil, false, err
	}

	switch typ {
	case model.ArtistType:
		artist := dbRes.(*model.Artist)
		res, err := findForExistingArtist(ctx, artist, services, repo, logger)
		return res, res.HasResults(), err

	case model.AlbumType:
		album := dbRes.(*model.Album)
		res, err := findForExistingAlbum(ctx, album, services, repo, logger)
		return res, res.HasResults(), err

	case model.TrackType:
		track := dbRes.(*model.Track)
		res, err := findForExistingTrack(ctx, track, services, repo, logger)
		return res, res.HasResults(), err

	case model.UnknownType:
		res, found, err := findNewThing(ctx, link, services, repo, logger)
		return res, found, err

	default:
		return nil, false, fmt.Errorf("unknown type %s", typ)
	}
}

func findForExistingArtist(ctx context.Context, foundArtist *model.Artist, services streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (*Result[*model.Artist], error) {

	logger = logger.WithField("artist_id", foundArtist.ArtistId)
	logger.Debugln("found an artist")

	res := NewResult[*model.Artist](model.ArtistType)

	// Find any related artists based on our artist ID
	existingArtists, err := repo.GetArtistsById(ctx, foundArtist.ArtistId)
	if err != nil {
		return nil, err
	}

	res.AddAll(existingArtists)

	// If we have results for all known services, then we're good to go
	if len(res.Items) == len(services) {
		return res, nil
	}

	logger.Debugf("looks like we have some new services since we found this artist (found %d, looking for %d)\n", len(existingArtists), len(services))

	// Query the remaining streaming service
	var newArtists []*model.Artist
	for key, service := range services {
		if res.HasResultFor(key) {
			continue
		}

		logger.Debugf("searching %s for artist\n", key)
		artist, found, err := service.SearchArtist(foundArtist)
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

	res.AddAll(newArtists)
	return res, nil
}

func findForExistingAlbum(ctx context.Context, foundAlbum *model.Album, services streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (*Result[*model.Album], error) {

	logger = logger.WithField("album_id", foundAlbum.AlbumId)
	logger.Debugln("found an album")

	res := NewResult[*model.Album](model.AlbumType)

	// Find any related albums based on our album ID
	existingAlbums, err := repo.GetAlbumsById(ctx, foundAlbum.AlbumId)
	if err != nil {
		return nil, err
	}

	res.AddAll(existingAlbums)

	// If we have results for all known services, then we're good to go
	if len(res.Items) == len(services) {
		return res, nil
	}

	logger.Debugf("looks like we have some new services since we found this album (found %d, looking for %d)\n", len(existingAlbums), len(services))

	// Query the remaining streaming service
	var newAlbums []*model.Album
	for key, service := range services {
		if res.HasResultFor(key) {
			continue
		}

		logger.Debugf("searching %s for album\n", key)
		album, found, err := service.SearchAlbum(foundAlbum)
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

	res.AddAll(newAlbums)
	return res, nil
}

func findForExistingTrack(ctx context.Context, foundTrack *model.Track, services streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (*Result[*model.Track], error) {

	logger = logger.WithField("isrc", foundTrack.Isrc)
	logger.Debugln("found a track")

	res := NewResult[*model.Track](model.TrackType)

	// Find any related track based on the ISRC
	existingTracks, err := repo.GetTracksByIsrc(ctx, foundTrack.Isrc)
	if err != nil {
		return nil, err
	}

	res.AddAll(existingTracks)

	// If we have results for all known services, then we're good to go
	if len(res.Items) == len(services) {
		return res, nil
	}

	logger.Debugf("looks like we have some new services since we found this track (found %d, looking for %d)\n", len(existingTracks), len(services))

	// Query the remaining streaming service
	var newTracks []*model.Track
	for key, service := range services {
		if res.HasResultFor(key) {
			continue
		}

		logger.Debugf("searching %s for track\n", key)
		track, found, err := service.GetTrackByIsrc(foundTrack.Isrc)
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

	res.AddAll(newTracks)
	return res, nil
}

func handleNewArtist(ctx context.Context, newArtist *model.Artist, services streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (*Result[*model.Artist], error) {

	res := NewResult[*model.Artist](model.ArtistType)
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
	for key, service := range services {
		if res.HasResultFor(key) {
			continue
		}

		logger.Debugf("searching %s for artist with name %s\n", key, newArtist.Name)
		foundArtist, found, err := service.SearchArtist(newArtist)
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

func handleNewAlbum(ctx context.Context, newAlbum *model.Album, services streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (*Result[*model.Album], error) {

	res := NewResult[*model.Album](model.AlbumType)
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
	for key, service := range services {
		if res.HasResultFor(key) {
			continue
		}

		logger.Debugf("searching %s for album with name %s\n", key, newAlbum.Name)
		foundAlbum, found, err := service.SearchAlbum(newAlbum)
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

func handleNewTrack(ctx context.Context, newTrack *model.Track, services streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (*Result[*model.Track], error) {

	res := NewResult[*model.Track](model.TrackType)
	res.Add(newTrack)

	logger = logger.WithField("isrc", newTrack.Isrc)

	newTracks := []*model.Track{
		newTrack,
	}

	// Query the other streaming services using what we found from the target streaming service
	for key, service := range services {
		logger.Debugf("searching %s for track with name %s\n", key, newTrack.Name)
		foundTrack, found, err := service.GetTrackByIsrc(newTrack.Isrc)
		if err != nil {
			logger.Errorf("%s: %s", service, err.Error())
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

func findNewThing(ctx context.Context, link string, services streamingservice.StreamingServices, repo db.Repository, logger *logrus.Entry) (any, bool, error) {

	logger.Debugln("looks like this is a new thing")

	// No links found, query the streaming service and find the same entry on other services
	var targetKey model.StreamingServiceType
	var targetService streamingservice.StreamingService
	otherServices := make(streamingservice.StreamingServices)

	// Figure out which streaming service the link belongs to
	for key, service := range services {
		if service.LinkBelongsToService(link) {
			targetKey = key
			targetService = service
		} else {
			otherServices[key] = service
		}
	}

	if targetService == nil {
		return nil, false, fmt.Errorf("couldn't find a streaming service for the given link: %s", link)
	}

	// Query the target streaming service
	logger.Debugf("searching %s\n", targetKey)
	typ, res, err := targetService.GetFromLink(link)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %s", targetKey, err.Error())
	}

	switch typ {
	case model.ArtistType:
		artist := res.(*model.Artist)
		res, err := handleNewArtist(ctx, artist, otherServices, repo, logger)
		return res, res.HasResults(), err

	case model.AlbumType:
		album := res.(*model.Album)
		res, err := handleNewAlbum(ctx, album, otherServices, repo, logger)
		return res, res.HasResults(), err

	case model.TrackType:
		track := res.(*model.Track)
		res, err := handleNewTrack(ctx, track, otherServices, repo, logger)
		return res, res.HasResults(), err

	case model.UnknownType:
		return nil, false, fmt.Errorf("could not find anything from %s", targetKey)

	default:
		return nil, false, fmt.Errorf("unknown type %s", typ)
	}
}
