package api

import (
	"context"
	"fmt"
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
	cfg *Config
	svr *http.Server
}

func NewMaestroApi(cfg *Config, svcCfg []streamingService.Config) (*MaestroApi, error) {

	cb := camogo.NewBuilder()
	if err := setupContainer(cb, svcCfg); err != nil {
		return nil, err
	}

	container := cb.Build()

	r := setupHandlers(container)

	addr := fmt.Sprintf(":%d", cfg.Port)
	svr := &http.Server{
		Addr: addr,

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,

		Handler:      r,
	}

	return &MaestroApi{cfg, svr}, nil
}


func (api *MaestroApi) Start() error {
	return api.svr.ListenAndServe()
}

func (api *MaestroApi) StartTLS() error {
	return api.svr.ListenAndServeTLS(api.cfg.CertFile, api.cfg.KeyFile)
}

func (api *MaestroApi) Shutdown() error {
	return api.svr.Shutdown(context.TODO())
}

func setupContainer(cb camogo.ContainerBuilder, svcCfg []streamingService.Config) error {

	// Todo: Use keys from config
	// Todo: Camogo needs slice support

	var services []streamingService.StreamingService
	for _, config := range svcCfg {
		if !config.Enabled {
			continue
		}

		switch config.ServiceName {
		case appleMusic.ConfigKey:
			token := config.Properties[appleMusic.TokenKey]
			services = append(services, appleMusic.NewAppleMusicStreamingService(token))

		case deezer.ConfigKey:
			services = append(services, deezer.NewDeezerStreamingService())

		case spotify.ConfigKey:
			clientId := config.Properties[spotify.ClientIdKey]
			clientSecret := config.Properties[spotify.ClientSecretKey]
			spService, err := spotify.NewSpotifyStreamingService(clientId, clientSecret)
			if err != nil {
				return err
			}

			services = append(services, spService)
		}
	}

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