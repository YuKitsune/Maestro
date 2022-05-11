package api

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yukitsune/maestro/pkg/api/apiconfig"
	mcontext "github.com/yukitsune/maestro/pkg/api/context"
	"github.com/yukitsune/maestro/pkg/api/db"
	"github.com/yukitsune/maestro/pkg/api/handlers"
	"github.com/yukitsune/maestro/pkg/api/middleware"
	"github.com/yukitsune/maestro/pkg/log"
	"github.com/yukitsune/maestro/pkg/metrics"
	"github.com/yukitsune/maestro/pkg/streamingservice"
	"github.com/yukitsune/maestro/pkg/streamingservice/applemusic"
	"github.com/yukitsune/maestro/pkg/streamingservice/deezer"
	"github.com/yukitsune/maestro/pkg/streamingservice/spotify"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"time"
)

type MaestroAPI struct {
	cfg    *apiconfig.Config
	logger *logrus.Entry
	svr    *http.Server
}

func NewMaestroAPI(apiCfg *apiconfig.Config, lCfg *log.Config, dbCfg *db.Config, svcCfg streamingservice.Config) (*MaestroAPI, error) {

	logger, err := configureLogger(context.Background(), lCfg)
	if err != nil {
		return nil, err
	}

	rec, err := configureMetrics()
	if err != nil {
		return nil, err
	}

	repo, err := configureRepository(context.Background(), dbCfg, rec, logger)
	if err != nil {
		return nil, err
	}

	services, err := configureServices(svcCfg, rec)
	if err != nil {
		return nil, err
	}

	router := configureHandlers(svcCfg, apiCfg, services, repo, rec, logger)

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

	api.logger.Warnln("TLS not enabled")
	api.logger.Infof("listening on %s", api.svr.Addr)
	return api.svr.ListenAndServe()
}

func (api *MaestroAPI) StartTLS() error {
	api.logger.Infoln("starting maestro API")
	api.logger.Debugf("config: %+v\n", viper.AllSettings())

	api.logger.Infof("listening on %s", api.svr.Addr)
	return api.svr.ListenAndServeTLS(api.cfg.CertFile, api.cfg.KeyFile)
}

func (api *MaestroAPI) Shutdown(ctx context.Context) error {
	api.logger.Infof("shutting down")
	return api.svr.Shutdown(ctx)
}

func configureLogger(ctx context.Context, cfg *log.Config) (*logrus.Entry, error) {

	logger := logrus.New()

	lvl, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	logger.SetLevel(lvl)

	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		PadLevelText:  true,
		FullTimestamp: true,
	})

	entry := logger.WithContext(ctx)

	reqID, err := mcontext.RequestID(ctx)
	if err == nil {
		entry = entry.WithField(log.RequestIDField, reqID)
	}

	return entry, nil
}

func configureMetrics() (metrics.Recorder, error) {
	return metrics.NewPrometheusMetricsRecorder()
}

func configureRepository(ctx context.Context, cfg *db.Config, rec metrics.Recorder, logger *logrus.Entry) (db.Repository, error) {
	opts := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	mdb := client.Database(cfg.Database)
	repo := db.NewMongoRepository(mdb, rec, logger)
	return repo, nil
}

func configureServices(svcCfg streamingservice.Config, rec metrics.Recorder) (streamingservice.StreamingServices, error) {

	var services []streamingservice.StreamingService
	for key, config := range svcCfg {
		if !config.Enabled() {
			continue
		}

		switch key {
		case applemusic.Key:
			cfg := config.(*applemusic.Config)
			s := applemusic.NewAppleMusicStreamingService(cfg, rec)
			services = append(services, s)
			break

		case deezer.Key:
			s := deezer.NewDeezerStreamingService(rec)
			services = append(services, s)
			break

		case spotify.Key:
			cfg := config.(*spotify.Config)
			s, err := spotify.NewSpotifyStreamingService(cfg, rec)
			if err != nil {
				return services, fmt.Errorf("failed to initialize spotify streaming service: %s", err.Error())
			}

			services = append(services, s)
			break
		}
	}

	return services, nil
}

func configureHandlers(scfg streamingservice.Config, acfg *apiconfig.Config, services streamingservice.StreamingServices, repo db.Repository, rec metrics.Recorder, logger *logrus.Entry) *mux.Router {

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
	r.HandleFunc("/services/{serviceName}/logo", handlers.GetServiceLogoHandler(scfg, acfg, logger)).Methods("GET")
	r.HandleFunc("/services", handlers.GetListServicesHandler(scfg)).Methods("GET")

	// Links
	r.HandleFunc("/link", handlers.GetLinkHandler(services, repo, logger)).Methods("GET").Queries("link", "{link}")
	r.HandleFunc("/artist/{id}", handlers.GetArtistByIdHandler(repo)).Methods("GET")
	r.HandleFunc("/album/{id}", handlers.GetAlbumByIdHandler(repo)).Methods("GET")
	r.HandleFunc("/track/{isrc}", handlers.GetTrackByIsrcHandler(repo, services, logger)).Methods("GET")

	return r
}
