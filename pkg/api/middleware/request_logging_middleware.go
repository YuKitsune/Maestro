package middleware

import (
	"github.com/sirupsen/logrus"
	"maestro/pkg/api/context"
	"maestro/pkg/api/handlers"
	"maestro/pkg/log"
	"net/http"
	"time"
)

func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqId, err := context.RequestId(r.Context())
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
			reqLogger := logger.WithField(log.RequestIdField, reqId).
				WithField("method", r.Method).
				WithField("path", r.URL.Path).
				WithField("status", rwd.StatusCode).
				WithField("duration", duration)

			q := r.URL.Query()
			if len(q) > 0 {
				reqLogger = reqLogger.WithField("query", q)
			}

			level := logrus.InfoLevel
			if isErrorCode(rwd.StatusCode) {
				level = logrus.ErrorLevel
			}

			reqLogger.Logln(level)
		})
	})
}
