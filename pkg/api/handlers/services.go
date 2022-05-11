package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/maestro/pkg/api/apiconfig"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingservice"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func GetListServicesHandler(cfg streamingservice.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var services []model.StreamingService
		for k, c := range cfg {
			sr := model.StreamingService{
				Name:    c.Name(),
				Key:     k.String(),
				Enabled: c.Enabled(),
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

func GetServiceLogoHandler(scfg streamingservice.Config, acfg *apiconfig.Config, logger *logrus.Entry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName, ok := vars["serviceName"]
		if !ok {
			BadRequest(w, "missing parameter \"serviceName\"")
			return
		}

		for k, c := range scfg {
			if k != model.StreamingServiceKey(serviceName) {
				continue
			}

			logger.Debugf("logo file name: %s", c.LogoFileName())

			logoFilePath := filepath.Join(acfg.AssetsDirectory, "logos", c.LogoFileName())
			logger.Debugf("logo path: %s", logoFilePath)

			// Ensure the file exists
			_, err := os.Stat(logoFilePath)
			if err != nil {
				exists := !errors.Is(err, os.ErrNotExist)
				if !exists {
					logger.Debugln("logo does not exist")
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
			return
		}

		NotFoundf(w, "couldn't find logo for %s", serviceName)
	}
}
