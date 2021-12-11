package handlers

import (
	mcontext "maestro/pkg/api/context"
	"maestro/pkg/streamingService"
	"net/http"
)

type serviceResource struct {
	Name        string
	Key         string
	ArtworkLink string
	Enabled     bool
}

func ListServices(w http.ResponseWriter, r *http.Request) {

	container, err := mcontext.Container(r.Context())
	if err != nil {
		Error(w, err)
		return
	}

	res, err := container.ResolveWithResult(func(cfg streamingService.Config) ([]serviceResource, error) {

		var res []serviceResource
		for k, c := range cfg {
			sr := serviceResource{
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

	Response(w, res, http.StatusOK)
}
