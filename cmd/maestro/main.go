package main

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yukitsune/maestro"
	"github.com/yukitsune/maestro/internal/grace"
	"github.com/yukitsune/maestro/pkg/api"
	"github.com/yukitsune/maestro/pkg/api/apiConfig"
	"github.com/yukitsune/maestro/pkg/api/db"
	"github.com/yukitsune/maestro/pkg/log"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingService"
	"github.com/yukitsune/maestro/pkg/streamingService/appleMusic"
	"github.com/yukitsune/maestro/pkg/streamingService/deezer"
	"github.com/yukitsune/maestro/pkg/streamingService/spotify"
	"strings"
	"time"
)

type Config struct {
	Api      *apiConfig.Config `mapstructure:"api"`
	Log      *log.Config       `mapstructure:"logging"`
	Db       *db.Config        `mapstructure:"database"`
	Services map[string]interface {
	} `mapstructure:"services"`
}

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
		grace.ExitFromError()
	}
}

func serve(_ *cobra.Command, _ []string) error {

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
			return err
		}
	}

	viper.AutomaticEnv()

	// Unmarshal
	var cfg *Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	scfg, err := decodeServiceConfigs(cfg)
	if err != nil {
		return err
	}

	maestroApi, err := api.NewMaestroApi(cfg.Api, cfg.Log, cfg.Db, scfg)
	if err != nil {
		grace.ExitFromError()
	}

	// Run our server in a goroutine so that it doesn't block.
	errorChan := make(chan error, 1)
	go func() {

		// Todo: TLS
		if err = maestroApi.Start(); err != nil {
			errorChan <- err
		}
	}()

	grace.WaitForShutdownSignalOrError(errorChan, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = maestroApi.Shutdown(ctx)
	})

	return nil
}

// Todo: It'd be nice to not have this...
func decodeServiceConfigs(c *Config) (streamingService.Config, error) {

	m := make(streamingService.Config)
	for k, v := range c.Services {

		key := model.StreamingServiceKey(k)
		var cfg streamingService.ServiceConfig

		switch key {
		case appleMusic.Key:
			var amc *appleMusic.Config
			if err := mapstructure.Decode(v, &amc); err != nil {
				return nil, err
			}

			cfg = amc
			break

		case deezer.Key:
			var dzc *deezer.Config
			if err := mapstructure.Decode(v, &dzc); err != nil {
				return nil, err
			}

			cfg = dzc
			break

		case spotify.Key:
			var spc *spotify.Config
			if err := mapstructure.Decode(v, &spc); err != nil {
				return nil, err
			}

			cfg = spc
			break

		default:
			return nil, fmt.Errorf("unknown service type %s", k)
		}

		m[key] = cfg
	}

	return m, nil
}
