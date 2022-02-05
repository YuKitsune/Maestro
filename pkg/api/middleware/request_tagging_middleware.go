package middleware

import (
	"github.com/google/uuid"
	"github.com/yukitsune/maestro/pkg/api/context"
	"net/http"
)

const RequestIDHeaderKey = "X-Request-Id"

func RequestTagging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()

		w.Header().Set(RequestIDHeaderKey, reqID)
		ctx := context.WithRequestID(r.Context(), reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
