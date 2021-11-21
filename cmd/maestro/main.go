package main

import (
	"maestro/internal/grace"
	"maestro/pkg/api"
)

func main() {

	maestroApi, err := api.NewMaestroApi()
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

	grace.WaitForShutdownSignalOrError(errorChan, func() { _ = maestroApi.Shutdown() })
}
