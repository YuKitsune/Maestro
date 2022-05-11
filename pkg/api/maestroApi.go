package api

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yukitsune/maestro/pkg/api/apiconfig"
	"github.com/yukitsune/maestro/pkg/api/db"
	"github.com/yukitsune/maestro/pkg/api/handlers"
	"github.com/yukitsune/maestro/pkg/api/middleware"
	"github.com/yukitsune/maestro/pkg/metrics"
	"github.com/yukitsune/maestro/pkg/streamingservice"
	"net/http"
	"time"
)

type MaestroAPI struct {
	cfg    *apiconfig.Config
	logger *logrus.Entry
	svr    *http.Server
}

func NewMaestroAPI(apiCfg *apiconfig.Config, serviceProvider streamingservice.ServiceProvider, repo db.Repository, rec metrics.Recorder, logger *logrus.Entry) (*MaestroAPI, error) {

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

	return &MaestroAPI{apiCfg, logger, svr}, nil
}

func (api *MaestroAPI) Start() error {
	api.logger.Infoln("starting maestro API")
	api.logger.Debugf("config: %+v\n", viper.AllSettings())

	api.logger.Infof("listening on %s", api.svr.Addr)
	return api.svr.ListenAndServe()
}

func (api *MaestroAPI) Shutdown(ctx context.Context) error {
	api.logger.Infof("shutting down")
	return api.svr.Shutdown(ctx)
}

func configureHandlers(cfg *apiconfig.Config, serviceProvider streamingservice.ServiceProvider, repo db.Repository, rec metrics.Recorder, logger *logrus.Entry) *mux.Router {

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
	r.HandleFunc("/services/{serviceName}/logo", handlers.GetServiceLogoHandler(cfg, serviceProvider, logger)).Methods("GET")
	r.HandleFunc("/services", handlers.GetListServicesHandler(serviceProvider)).Methods("GET")

	// Links
	r.HandleFunc("/link", handlers.GetLinkHandler(serviceProvider, repo, logger)).Methods("GET").Queries("link", "{link}")
	r.HandleFunc("/artist/{id}", handlers.GetArtistByIdHandler(repo)).Methods("GET")
	r.HandleFunc("/album/{id}", handlers.GetAlbumByIdHandler(repo)).Methods("GET")
	r.HandleFunc("/track/{isrc}", handlers.GetTrackByIsrcHandler(repo, serviceProvider, logger)).Methods("GET")

	return r
}
