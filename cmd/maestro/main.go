package main

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	"github.com/yukitsune/maestro/pkg/streamingservice"
	"github.com/yukitsune/maestro/pkg/streamingservice/provider"
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

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	logger, err := configureLogger(cfg.Log)
	if err != nil {
		return err
	}

	logger.Debugf("Config: %+v", cfg)

	rec, err := configureMetrics()
	if err != nil {
		return err
	}

	repo, err := configureRepository(cfg.Database, rec, logger)
	if err != nil {
		return err
	}

	serviceProvider, err := configureServices(cfg.Services, rec)
	if err != nil {
		return err
	}

	maestroAPI, err := api.NewMaestroServer(cfg.API, serviceProvider, repo, rec, logger)
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

func loadConfig() (*config.Config, error) {

	// Environment variables
	viper.SetEnvPrefix("MAESTRO")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// Config file
	viper.SetConfigName("maestro")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/maestro")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Load in the configuration
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
			fmt.Println("Config file not found")
		} else {
			return nil, err
		}
	}

	viper.AutomaticEnv()

	// Unmarshal
	var cfg *config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func configureLogger(cfg *config.Log) (*logrus.Logger, error) {

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

	if cfg.Loki != nil {

		// Grafana doesn't have a "panic" level, but it does have a "critical" level
		// https://grafana.com/docs/grafana/latest/explore/logs-integration/
		opts := lokirus.NewLokiHookOptions().
			WithLevelMap(lokirus.LevelMap{logrus.PanicLevel: "critical"}).
			WithStaticLabels(cfg.Loki.Labels)
		hook := lokirus.NewLokiHookWithOpts(
			cfg.Loki.Host,
			opts,
			logrus.InfoLevel,
			logrus.WarnLevel,
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel)

		logger.AddHook(hook)
	}

	return logger, nil
}

func configureMetrics() (metrics.Recorder, error) {
	return metrics.NewPrometheusMetricsRecorder()
}

func configureRepository(cfg *config.Database, rec metrics.Recorder, logger *logrus.Logger) (db.Repository, error) {
	opts := options.Client().ApplyURI(cfg.URI)
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

	mdb := client.Database(cfg.Database)
	repo := db.NewMongoRepository(mdb, rec, logger)
	return repo, nil
}

func configureServices(cfg *config.Services, rec metrics.Recorder) (streamingservice.ServiceProvider, error) {
	return provider.NewDefaultProvider(cfg, rec)
}
