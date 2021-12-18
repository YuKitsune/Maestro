package middleware

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"maestro/pkg/api/context"
	"maestro/pkg/api/handlers"
	"net/http"
	"runtime"
)

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				// Todo: It'd be nice to not have to fetch the container here
				container, _ := context.Container(r.Context())
				if container != nil {
					_ = container.Resolve(func(logger *logrus.Entry) {
						logger.
							WithField("stacktrace", fmt.Sprintf("%s", buf)).
							Errorf("recovering from panic: %s\n", err)
					})
				}

				// Write error message only
				handlers.Error(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
