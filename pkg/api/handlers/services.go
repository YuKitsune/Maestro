package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/maestro/pkg/api/apiConfig"
	mcontext "github.com/yukitsune/maestro/pkg/api/context"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingService"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func ListServices(w http.ResponseWriter, r *http.Request) {

	container, err := mcontext.Container(r.Context())
	if err != nil {
		Error(w, err)
		return
	}

	res, err := container.ResolveWithResult(func(cfg streamingService.Config) ([]model.StreamingService, error) {

		var res []model.StreamingService
		for k, c := range cfg {
			sr := model.StreamingService{
				Name:    c.Name(),
				Key:     k.String(),
				Enabled: c.Enabled(),
			}

			res = append(res, sr)
		}

		return res, nil
	})

	if err != nil {
		Error(w, err)
		return
	}

	services := res.([]model.StreamingService)
	if res == nil || len(services) == 0 {
		NotFound(w, "could not find any services")
		return
	}

	Response(w, res, http.StatusOK)
}

func GetLogo(w http.ResponseWriter, r *http.Request) {

	container, err := mcontext.Container(r.Context())
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

	res, err := container.ResolveWithResult(func(scfg streamingService.Config, acfg *apiConfig.Config, logger *logrus.Entry) ([]byte, error) {
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
					return []byte{}, nil
				} else {
					return nil, err
				}
			}

			logo, err := ioutil.ReadFile(logoFilePath)
			return logo, err
		}

		return nil, nil
	})

	if err != nil {
		Error(w, err)
		return
	}

	if res == nil {
		NotFoundf(w, "could not find service with name %s", serviceName)
		return
	}

	bytes := res.([]byte)
	if len(bytes) <= 0 {
		NotFoundf(w, "could not find logo for service %s", serviceName)
		return
	}

	Image(w, bytes)
}
