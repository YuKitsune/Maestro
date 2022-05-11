package middleware

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/maestro/pkg/api/handlers"
	"net/http"
	"runtime"
)

func PanicRecovery(logger *logrus.Entry) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					buf := make([]byte, 2048)
					n := runtime.Stack(buf, false)
					buf = buf[:n]

					logger.WithField("stacktrace", fmt.Sprintf("%s", buf)).
						Errorf("recovering from panic: %s\n", err)

					// Write error message only
					handlers.Error(w, fmt.Errorf("%s", err))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}

}
