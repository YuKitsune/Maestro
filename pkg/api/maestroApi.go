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

func NewMaestroApi(apiCfg *Config, dbCfg *db.Config, scfg streamingService.Config) (*MaestroApi, error) {

	cb := camogo.NewBuilder()
	if err := setupContainer(cb, dbCfg, scfg); err != nil {
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

func (api *MaestroApi) Shutdown(ctx context.Context) error {
	return api.svr.Shutdown(ctx)
}

func setupContainer(cb camogo.ContainerBuilder, dbCfg *db.Config, scfg streamingService.Config) error {

	// Todo: Context timeout here
	if err := cb.RegisterFactory(func() context.Context {
		return context.TODO()
	}, camogo.SingletonLifetime); err != nil {
		return nil
	}

	dbMod := &db.DatabaseModule{Config: dbCfg}
	if err := cb.RegisterModule(dbMod); err != nil {
		return err
	}

	if err := registerStreamingServices(cb, scfg); err != nil {
		return err
	}

	return nil
}

func registerStreamingServices(cb camogo.ContainerBuilder, scfg streamingService.Config) error {

	if err := cb.RegisterInstance(scfg); err != nil {
		return err
	}

	// Todo: Camogo needs slice support
	factory := func (c streamingService.Config) ([]streamingService.StreamingService, error) {
		var services []streamingService.StreamingService
		for key, config := range c {
			if !config.Enabled() {
				continue
			}

			switch key {
			case appleMusic.Key:
				cfg := config.(*appleMusic.Config)
				s := appleMusic.NewAppleMusicStreamingService(cfg)
				services = append(services, s)
				break

			case deezer.Key:
				s := deezer.NewDeezerStreamingService()
				services = append(services, s)
				break

			case spotify.Key:
				cfg := config.(*spotify.Config)
				s, err := spotify.NewSpotifyStreamingService(cfg)
				if err != nil {
					return services, err
				}

				services = append(services, s)
				break
			}
		}

		return services, nil
	}

	// Need to register these as transient (scoped) because Spotify doesn't provide refresh tokens for client credentials
	if err := cb.RegisterFactory(factory, camogo.ScopedLifetime); err != nil {
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
	r.HandleFunc("/{groupId}", handlers.HandleGroup).Methods("GET")
	// r.HandleFunc("/flag", handlers.HandleFlag).Methods("GET").Queries("thingId", "{thingId}", "thingType", "{thingType}")

	return r
}
