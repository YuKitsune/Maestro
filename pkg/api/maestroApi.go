package api

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/yukitsune/camogo"
	"maestro/pkg/api/handlers"
	"maestro/pkg/api/middleware"
	"maestro/pkg/streamingService"
	"maestro/pkg/streamingService/appleMusic"
	"maestro/pkg/streamingService/deezer"
	"maestro/pkg/streamingService/spotify"
	"net/http"
	"time"
)

type MaestroApi struct {
	svr *http.Server
}

func NewMaestroApi() (*MaestroApi, error) {

	cb := camogo.NewBuilder()
	if err := setupContainer(cb); err != nil {
		return nil, err
	}

	container := cb.Build()

	r := setupHandlers(container)
	svr := &http.Server{
		// Todo: Implement once we have configs
		// Addr: api.config.GetAddress(),

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,

		Handler:      r,
	}

	return &MaestroApi{svr}, nil
}


func (api *MaestroApi) Start() error {
	return api.svr.ListenAndServe()
}

func (api *MaestroApi) StartTLS() error {

	// Todo: Implement once we have configs
	// return api.svr.ListenAndServeTLS(api.config.CertFile, api.config.KeyFile)

	return api.Start()
}

func (api *MaestroApi) Shutdown() error {
	return api.svr.Shutdown(context.TODO())
}

func setupContainer(cb camogo.ContainerBuilder) error {

	// Todo: Use keys from config
	// Todo: Camogo needs slice support
	services = append(services, appleMusic.NewAppleMusicStreamingService(amToken))
	services = append(services, spotify.NewSpotifyStreamingService(spToken))
	services = append(services, deezer.NewDeezerStreamingService())

	if err := cb.RegisterInstance(services); err != nil {
		return err
	}

	return nil
}

func setupHandlers(container camogo.Container) *mux.Router {

	r := mux.NewRouter()

	// Middleware
	containerInjectionMiddleware := middleware.NewContainerInjectionMiddleware(container)
	r.Use(containerInjectionMiddleware.Middleware)

	// Routes
	// Search
	r.HandleFunc("/search/artist", handlers.HandleSearchArtist).Methods("GET").Queries("name", "{name}")
	r.HandleFunc("/search/album", handlers.HandleSearchAlbum).Methods("GET").Queries("name", "{name}")
	r.HandleFunc("/search/song", handlers.HandleSearchSong).Methods("GET").Queries("name", "{name}")

	return r
}