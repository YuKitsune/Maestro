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

			// Get the route template
			route := mux.CurrentRoute(r)
			pathTemplate, _ := route.GetPathTemplate()

			// Record request duration with the path
			rec.ReportRequestDuration(pathTemplate, func() {
				next.ServeHTTP(&rwd, r)
			})

			if isServerErrorCode(rwd.StatusCode) {
				go rec.CountServerError()
			}
		})
	}
}
