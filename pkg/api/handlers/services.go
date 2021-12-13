package handlers

import (
	"fmt"
	mcontext "maestro/pkg/api/context"
	"maestro/pkg/model"
	"maestro/pkg/streamingService"
	"net/http"
)

func ListServices(w http.ResponseWriter, r *http.Request) {

	container, err := mcontext.Container(r.Context())
	if err != nil {
		Error(w, err)
		return
	}

	res, err := container.ResolveWithResult(func(cfg streamingService.Config) ([]model.StreamingService, error) {

		return nil, fmt.Errorf("fuck you")

		var res []model.StreamingService
		for k, c := range cfg {
			sr := model.StreamingService{
				Name:        c.Name(),
				Key:         k.String(),
				ArtworkLink: c.ArtworkLink(),
				Enabled:     c.Enabled(),
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
