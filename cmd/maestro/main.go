package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yukitsune/maestro/pkg/streamingservice/provider"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yukitsune/lokirus"
	"github.com/yukitsune/maestro"
	"github.com/yukitsune/maestro/internal/grace"
	"github.com/yukitsune/maestro/pkg/api"
	"github.com/yukitsune/maestro/pkg/api/db"
	"github.com/yukitsune/maestro/pkg/config"
	"github.com/yukitsune/maestro/pkg/metrics"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(maestro.Version)
			return nil
		},
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the REST API",
		RunE:  serve,
	}

	rootCmd := &cobra.Command{
		Use:   "maestro <command> [flags]",
		Short: "The Maestro REST API",
	}

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)

	err := rootCmd.Execute()
	if err != nil {
		grace.ExitFromError(err)
	}
}

func serve(_ *cobra.Command, _ []string) error {

	logger := logrus.New()
	v := viper.New()

	cfg := loadConfig(v, logger)

	configureLogger(cfg.Logging(), logger)

	// When using the debug log level, print the config out
	logger.Debugf("Config: %+v", cfg.Debug())

	rec, err := configureMetrics()
	if err != nil {
		return err
	}

	repo, err := configureRepository(cfg.Database(), rec, logger)
	if err != nil {
		return err
	}

	serviceProvider, err := provider.NewDefaultProvider(cfg.Services(), rec)
	if err != nil {
		return err
	}

	maestroAPI, err := api.NewMaestroServer(cfg.API(), serviceProvider, repo, rec, logger)
	if err != nil {
		grace.ExitFromError(err)
	}

	// Run our server in a goroutine so that it doesn't block.
	errorChan := make(chan error, 1)
	go func() {
		if err = maestroAPI.Start(); err != nil {
			errorChan <- err
		}
	}()

	grace.WaitForShutdownSignalOrError(errorChan, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = maestroAPI.Shutdown(ctx)
	})

	return nil
}

func loadConfig(v *viper.Viper, logger logrus.FieldLogger) config.Config {

	// Config file
	v.SetConfigName("maestro")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/maestro")
	v.AddConfigPath("../configs")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	// Watch for changes
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		logger.Infof("Config file changed: ", e.Name)
	})

	_ = v.ReadInConfig()

	// Environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.SetEnvPrefix("MAESTRO")
	v.AutomaticEnv()

	return config.NewViperConfig(v)
}

func configureLogger(cfg config.Logging, logger *logrus.Logger) {

	logger.SetLevel(cfg.Level())

	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		PadLevelText:  true,
		FullTimestamp: true,
	})

	if cfg.Loki().Enabled() {

		// Grafana doesn't have a "panic" level, but it does have a "critical" level
		// https://grafana.com/docs/grafana/latest/explore/logs-integration/
		opts := lokirus.NewLokiHookOptions().
			WithLevelMap(lokirus.LevelMap{logrus.PanicLevel: "critical"}).
			WithStaticLabels(cfg.Loki().Labels())
		hook := lokirus.NewLokiHookWithOpts(
			cfg.Loki().Host(),
			opts,
			logrus.InfoLevel,
			logrus.WarnLevel,
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel)

		logger.AddHook(hook)
	}
}

func configureMetrics() (metrics.Recorder, error) {
	return metrics.NewPrometheusMetricsRecorder()
}

func configureRepository(cfg config.Database, rec metrics.Recorder, logger *logrus.Logger) (db.Repository, error) {
	opts := options.Client().ApplyURI(cfg.Uri())
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	mdb := client.Database(cfg.Name())
	repo := db.NewMongoRepository(mdb, rec, logger)
	return repo, nil
}
