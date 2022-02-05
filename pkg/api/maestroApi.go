package api

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yukitsune/camogo"
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
	"net/http"
	"time"
)

type MaestroAPI struct {
	cfg    *apiconfig.Config
	logger *logrus.Entry
	svr    *http.Server
}

func NewMaestroAPI(apiCfg *apiconfig.Config, lCfg *log.Config, dbCfg *db.Config, scfg streamingservice.Config) (*MaestroAPI, error) {

	logger, err := configureLogger(context.Background(), lCfg)
	if err != nil {
		return nil, err
	}

	cb := camogo.NewBuilder()
	if err := setupContainer(cb, apiCfg, lCfg, dbCfg, scfg); err != nil {
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

func setupContainer(cb camogo.ContainerBuilder, aCfg *apiconfig.Config, lCfg *log.Config, dbCfg *db.Config, sCfg streamingservice.Config) error {

	// Todo: Context timeout here
	if err := cb.RegisterFactory(func() context.Context {
		return context.TODO()
	}, camogo.SingletonLifetime); err != nil {
		return nil
	}

	// Log Config
	if err := cb.RegisterInstance(lCfg); err != nil {
		return err
	}

	// API Config
	if err := cb.RegisterInstance(aCfg); err != nil {
		return err
	}

	// Logger
	if err := cb.RegisterFactory(configureLogger, camogo.ScopedLifetime); err != nil {
		return err
	}

	// Metrics
	if err := cb.RegisterFactory(metrics.NewPrometheusMetricsRecorder, camogo.SingletonLifetime); err != nil {
		return err
	}

	// Database
	dbMod := &db.DatabaseModule{Config: dbCfg}
	if err := cb.RegisterModule(dbMod); err != nil {
		return err
	}

	// Streaming services
	if err := registerStreamingServices(cb, sCfg); err != nil {
		return err
	}

	return nil
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

func registerStreamingServices(cb camogo.ContainerBuilder, scfg streamingservice.Config) error {

	if err := cb.RegisterInstance(scfg); err != nil {
		return err
	}

	// Todo: Camogo needs slice support
	factory := func(c streamingservice.Config, mr metrics.Recorder) ([]streamingservice.StreamingService, error) {

		var services []streamingservice.StreamingService
		for key, config := range c {
			if !config.Enabled() {
				continue
			}

			switch key {
			case applemusic.Key:
				cfg := config.(*applemusic.Config)
				s := applemusic.NewAppleMusicStreamingService(cfg, mr)
				services = append(services, s)
				break

			case deezer.Key:
				s := deezer.NewDeezerStreamingService(mr)
				services = append(services, s)
				break

			case spotify.Key:
				cfg := config.(*spotify.Config)
				s, err := spotify.NewSpotifyStreamingService(cfg, mr)
				if err != nil {
					return services, fmt.Errorf("failed to initialize spotify streaming service: %s", err.Error())
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
	r.Use(middleware.RequestTagging)

	containerInjectionMiddleware := middleware.NewContainerInjectionMiddleware(container)
	r.Use(containerInjectionMiddleware.Middleware)

	r.Use(middleware.Metrics)
	r.Use(middleware.RequestLogging)
	r.Use(middleware.PanicRecovery)
	r.Use(middleware.Cors)

	// Routes
	r.NotFoundHandler = http.HandlerFunc(handlers.HandleNotFound)

	// Metrics
	r.Handle("/metrics", promhttp.Handler())

	// Services
	r.HandleFunc("/services/{serviceName}/logo", handlers.GetLogo).Methods("GET")
	r.HandleFunc("/services", handlers.ListServices).Methods("GET")

	// Links
	r.HandleFunc("/link", handlers.HandleLink).Methods("GET").Queries("link", "{link}")
	r.HandleFunc("/{groupId}", handlers.HandleGroup).Methods("GET")

	return r
}
