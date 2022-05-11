package grace

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type shutdownHook func()

func WaitForShutdownSignal(shutdownHooks ...shutdownHook) {
	waitForSignal(getShutdownSignalChan(), make(chan error, 1))
	for _, hook := range shutdownHooks {
		hook()
	}
}

func WaitForShutdownSignalOrError(errorChan chan error, shutdownHooks ...shutdownHook) {
	waitForSignal(getShutdownSignalChan(), errorChan)
	for _, hook := range shutdownHooks {
		hook()
	}
}

func waitForSignal(shutdownSignalChan chan os.Signal, errorChan chan error) {
	select {
	case sig := <-shutdownSignalChan:
		handleShutdownSignal(sig)
		break

	case err := <-errorChan:
		ExitFromError(err)
		break
	}
}

func handleShutdownSignal(sig os.Signal) {
	if sig == syscall.SIGTERM || sig == syscall.SIGQUIT || sig == syscall.SIGINT || sig == os.Kill {
		log.Println("Shutdown signal caught")
		go func() {
			select {
			// exit if graceful shutdown not finished in 60 sec.
			case <-time.After(time.Second * 60):
				log.Fatalln("graceful shutdown timed out")
			}
		}()
		log.Println("Shutdown completed, exiting")
		return
	}

	log.Printf("Shutdown, unknown signal caught: %s", sig.String())
	return
}

func getShutdownSignalChan() chan os.Signal {
	shutdownSignalChan := make(chan os.Signal, 1)
	signal.Notify(shutdownSignalChan,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGKILL,
		os.Kill,
	)

	return shutdownSignalChan
}

func ExitFromError(err error) {
	log.Fatalf("an error has occurred: %s", err.Error())
}
