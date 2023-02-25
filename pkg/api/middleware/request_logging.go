package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/maestro/pkg/api/context"
	"github.com/yukitsune/maestro/pkg/api/responses"
	"github.com/yukitsune/maestro/pkg/log"
)

func RequestLogging(logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			reqID, err := context.RequestID(r.Context())
			if err != nil {
				responses.Error(w, err)
				return
			}

			// wrap original http.ResponseWriter
			rwd := responseWriterDecorator{
				ResponseWriter: w,
				StatusCode:     0,
			}

			start := time.Now()
			next.ServeHTTP(&rwd, r)
			duration := time.Since(start)

			reqLogger := logger.WithField(log.RequestIDField, reqID).
				WithField("method", r.Method).
				WithField("path", r.URL.Path).
				WithField("status", rwd.StatusCode).
				WithField("duration", duration)

			q := r.URL.Query()
			if len(q) > 0 {
				reqLogger = reqLogger.WithField("query", q)
			}

			level := logrus.DebugLevel
			if isServerErrorCode(rwd.StatusCode) {
				level = logrus.ErrorLevel
			}

			reqLogger.Logln(level)
		})
	}
}
