package middleware

import (
	"fmt"
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

				// Todo: Log error
				// Todo: It'd be nice to not have to fetch the container here
				//container, _ := context.Container(r.Context())
				//if container != nil {
				//	_ = container.Resolve(func(logger *logrus.Logger) {
				//		logger.Errorf("Recovering from panic: %v\n%s\n", err, buf)
				//	})
				//}

				handlers.Error(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
