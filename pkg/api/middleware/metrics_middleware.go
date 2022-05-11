package middleware

import (
	"github.com/gorilla/mux"
	"github.com/yukitsune/maestro/pkg/metrics"
	"net/http"
)

func Metrics(rec metrics.Recorder) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Don't want metric requests messing with our actual metrics
			if r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r)
				return
			}

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
			}
		})
	}
}
