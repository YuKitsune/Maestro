package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yukitsune/maestro/pkg/api/db"
	"github.com/yukitsune/maestro/pkg/api/handlers"
	"github.com/yukitsune/maestro/pkg/api/middleware"
	"github.com/yukitsune/maestro/pkg/config"
	"github.com/yukitsune/maestro/pkg/metrics"
	"github.com/yukitsune/maestro/pkg/streamingservice"
)

type MaestroServer struct {
	cfg    *config.API
	logger *logrus.Logger
	svr    *http.Server
}

func NewMaestroServer(apiCfg *config.API, serviceProvider streamingservice.ServiceProvider, repo db.Repository, rec metrics.Recorder, logger *logrus.Logger) (*MaestroServer, error) {

	router := configureHandlers(apiCfg, serviceProvider, repo, rec, logger)

	addr := fmt.Sprintf(":%d", apiCfg.Port)
	svr := &http.Server{
		Addr: addr,

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,

		Handler: router,
	}

	return &MaestroServer{apiCfg, logger, svr}, nil
}

func (api *MaestroServer) Start() error {
	api.logger.Infoln("starting maestro server")
	api.logger.Debugf("config: %+v\n", viper.AllSettings())

	api.logger.Infof("listening on %s", api.svr.Addr)
	return api.svr.ListenAndServe()
}

func (api *MaestroServer) Shutdown(ctx context.Context) error {
	api.logger.Infof("shutting down maestro server")
	return api.svr.Shutdown(ctx)
}

func configureHandlers(apiConfig *config.API, serviceProvider streamingservice.ServiceProvider, repo db.Repository, rec metrics.Recorder, logger *logrus.Logger) *mux.Router {

	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.RequestTagging)

	r.Use(middleware.Metrics(rec))
	r.Use(middleware.RequestLogging(logger))
	r.Use(middleware.PanicRecovery(logger))
	r.Use(middleware.Cors)

	// Routes
	r.NotFoundHandler = http.HandlerFunc(handlers.HandleNotFound)

	// Metrics
	r.Handle("/metrics", promhttp.Handler())

	// Services
	r.HandleFunc("/services/{serviceName}/logo", handlers.GetServiceLogoHandler(apiConfig, serviceProvider, logger)).Methods("GET")
	r.HandleFunc("/services", handlers.GetListServicesHandler(serviceProvider)).Methods("GET")

	// Links
	r.HandleFunc("/link", handlers.GetLinkHandler(serviceProvider, repo, logger)).Methods("GET").Queries("link", "{link}")
	r.HandleFunc("/artist/{id}", handlers.GetArtistByIdHandler(repo)).Methods("GET")
	r.HandleFunc("/album/{id}", handlers.GetAlbumByIdHandler(repo)).Methods("GET")
	r.HandleFunc("/track/{isrc}", handlers.GetTrackByIsrcHandler(repo, serviceProvider, logger)).Methods("GET")

	return r
}
