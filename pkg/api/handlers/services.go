package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/maestro/pkg/api/apiconfig"
	"github.com/yukitsune/maestro/pkg/log"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingservice"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func GetListServicesHandler(sp streamingservice.ServiceProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svcs := sp.ListConfigs()

		var services []model.StreamingService
		for k, cfg := range svcs {
			sr := model.StreamingService{
				Name:    cfg.Name(),
				Key:     k.String(),
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

func GetServiceLogoHandler(apiCfg *apiconfig.Config, sp streamingservice.ServiceProvider, logger *logrus.Logger) http.HandlerFunc {
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

		cfg, err := sp.GetConfig(model.StreamingServiceKey(serviceName))
		if err != nil {
			NotFoundf(w, "couldn't find streaming service with key %s", serviceName)
			return
		}

		reqLogger.Debugf("logo file name: %s", cfg.LogoFileName())

		logoFilePath := filepath.Join(apiCfg.AssetsDirectory, "logos", cfg.LogoFileName())
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
