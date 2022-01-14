package middleware

import (
	"net/http"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		w.Header().Set("Access-Control-Allow-Origin", "*")
	})
}
