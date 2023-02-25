package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/maestro/pkg/api/responses"
	"github.com/yukitsune/maestro/pkg/log"
)

func PanicRecovery(logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					buf := make([]byte, 2048)
					n := runtime.Stack(buf, false)
					buf = buf[:n]

					reqLogger, logErr := log.ForRequest(logger, r)
					if logErr != nil {
						responses.Error(w, logErr)
					}

					reqLogger.WithField("stacktrace", fmt.Sprintf("%s", buf)).
						Errorf("recovering from panic: %s\n", err)

					// Write error message only
					responses.Error(w, fmt.Errorf("%s", err))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}

}
