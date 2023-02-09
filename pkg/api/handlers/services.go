package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/maestro/pkg/config"
	"github.com/yukitsune/maestro/pkg/log"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingservice"
)

func GetListServicesHandler(serviceProvider streamingservice.ServiceProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var services []model.StreamingService
		for _, cfg := range serviceProvider.ListConfigs() {
			sr := model.StreamingService{
				Key:     cfg.Type(),
				Name:    cfg.Name(),
				Enabled: cfg.Enabled(),
			}

			services = append(services, sr)
		}

		if len(services) == 0 {
			NotFound(w, "could not find any services")
			return
		}

		Response(w, services, http.StatusOK)
	}
}

func GetServiceLogoHandler(apiConfig config.API, serviceProvider streamingservice.ServiceProvider, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLogger, err := log.ForRequest(logger, r)
		if err != nil {
			Error(w, err)
			return
		}

		vars := mux.Vars(r)
		serviceName, ok := vars["serviceName"]
		if !ok {
			BadRequest(w, "missing parameter \"serviceName\"")
			return
		}

		cfg, err := serviceProvider.GetConfig(model.StreamingServiceType(serviceName))
		if err != nil {
			NotFoundf(w, "couldn't find streaming service with key %s", serviceName)
			return
		}

		reqLogger.Debugf("logo file name: %s", cfg.LogoFileName())

		logoFilePath := filepath.Join(apiConfig.AssetsDirectory(), "logos", cfg.LogoFileName())
		reqLogger.Debugf("logo path: %s", logoFilePath)

		// Ensure the file exists
		_, err = os.Stat(logoFilePath)
		if err != nil {
			exists := !errors.Is(err, os.ErrNotExist)
			if !exists {
				reqLogger.Debugln("logo does not exist")
				NotFoundf(w, "couldn't find logo for %s", serviceName)
				return
			}

			Error(w, err)
			return
		}

		logo, err := ioutil.ReadFile(logoFilePath)
		if err != nil {
			Error(w, err)
			return
		}

		Image(w, logo)
	}
}
