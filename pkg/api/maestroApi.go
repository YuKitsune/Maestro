package api

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/yukitsune/camogo"
	"maestro/pkg/api/db"
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

func NewMaestroApi(apiCfg *Config, dbCfg *db.Config, svgCfgs []streamingService.Config) (*MaestroApi, error) {

	cb := camogo.NewBuilder()
	if err := setupContainer(cb, dbCfg, svgCfgs); err != nil {
		return nil, err
	}

	container := cb.Build()

	r := setupHandlers(container)

	addr := fmt.Sprintf(":%d", apiCfg.Port)
	svr := &http.Server{
		Addr: addr,

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,

		Handler: r,
	}

	return &MaestroApi{apiCfg, svr}, nil
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

func setupContainer(cb camogo.ContainerBuilder, dbCfg *db.Config, svcCfgs []streamingService.Config) error {

	dbMod := &db.DatabaseModule{Config: dbCfg}
	if err := cb.RegisterModule(dbMod); err != nil {
		return err
	}

	if err := registerStreamingServices(cb, svcCfgs); err != nil {
		return err
	}

	return nil
}

func registerStreamingServices(cb camogo.ContainerBuilder, svcCfg []streamingService.Config) error {

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
			shareUrlPattern := config.Properties[streamingService.ShareLinkPatternKey]
			services = append(services, appleMusic.NewAppleMusicStreamingService(token, shareUrlPattern))

		case deezer.ConfigKey:
			shareUrlPattern := config.Properties[streamingService.ShareLinkPatternKey]
			services = append(services, deezer.NewDeezerStreamingService(shareUrlPattern))

		case spotify.ConfigKey:
			clientId := config.Properties[spotify.ClientIdKey]
			clientSecret := config.Properties[spotify.ClientSecretKey]
			shareUrlPattern := config.Properties[streamingService.ShareLinkPatternKey]
			spService, err := spotify.NewSpotifyStreamingService(clientId, clientSecret, shareUrlPattern)
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

	// Links
	r.HandleFunc("/link", handlers.HandleLink).Methods("GET").Queries("link", "{link}")
	r.HandleFunc("/link/flag", handlers.HandleFlagLink).Methods("GET").Queries("linkId", "{linkId}")

	return r
}
