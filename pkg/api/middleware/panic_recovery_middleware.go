package middleware

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/maestro/pkg/api/handlers"
	"github.com/yukitsune/maestro/pkg/log"
	"net/http"
	"runtime"
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
						handlers.Error(w, logErr)
					}

					reqLogger.WithField("stacktrace", fmt.Sprintf("%s", buf)).
						Errorf("recovering from panic: %s\n", err)

					// Write error message only
					handlers.Error(w, fmt.Errorf("%s", err))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}

}
