package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"maestro"
	"maestro/internal/grace"
	"maestro/pkg/api"
	"maestro/pkg/api/db"
	"maestro/pkg/streamingService"
	"time"
)

type Config struct {
	Api      *api.Config               `mapstructure:"api"`
	Db       *db.Config                `mapstructure:"database"`
	Services []streamingService.Config `mapstructure:"services"`
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
		grace.ExitFromError(err)
	}
}

func serve(_ *cobra.Command, _ []string) error {

	// Config file
	viper.SetConfigName("maestro")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/chameleon")
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

	// Unmarshal
	var cfg *Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	maestroApi, err := api.NewMaestroApi(cfg.Api, cfg.Db, cfg.Services)
	if err != nil {
		grace.ExitFromError(err)
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
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()
		_ = maestroApi.Shutdown(ctx)
	})

	return nil
}
