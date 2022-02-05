package middleware

import (
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/maestro/pkg/api/context"
	"github.com/yukitsune/maestro/pkg/api/handlers"
	"github.com/yukitsune/maestro/pkg/log"
	"net/http"
	"time"
)

func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqID, err := context.RequestID(r.Context())
		if err != nil {
			handlers.Error(w, err)
			return
		}

		ctr, err := context.Container(r.Context())
		if err != nil {
			handlers.Error(w, err)
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

		_ = ctr.Resolve(func(logger *logrus.Entry) {
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
			} else if isClientErrorCode(rwd.StatusCode) {
				level = logrus.WarnLevel
			}

			reqLogger.Logln(level)
		})
	})
}
