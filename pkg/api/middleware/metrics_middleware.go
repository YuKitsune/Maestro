package middleware

import (
	"maestro/pkg/api/context"
	"maestro/pkg/api/handlers"
	"maestro/pkg/metrics"
	"net/http"
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Don't want metric requests messing with our actual metrics
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		ctr, err := context.Container(r.Context())
		if err != nil {
			handlers.Error(w, err)
			return
		}

		err = ctr.Resolve(func(rec metrics.Recorder) error {

			// wrap original http.ResponseWriter
			rwd := responseWriterDecorator{
				ResponseWriter: w,
				StatusCode:     0,
			}

			rec.ReportRequestDuration(func() {
				next.ServeHTTP(&rwd, r)
			})

			go rec.CountRequest()

			if isServerErrorCode(rwd.StatusCode) {
				go rec.CountServerError()
			} else if isClientErrorCode(rwd.StatusCode) {
				go rec.CountClientError()
			}

			return nil
		})

		if err != nil {
			handlers.Error(w, err)
			return
		}
	})
}
