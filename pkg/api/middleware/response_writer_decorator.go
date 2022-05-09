package middleware

import "net/http"

// we use this in logging and metrics middleware
type responseWriterDecorator struct {
	http.ResponseWriter
	StatusCode int
}

func (r *responseWriterDecorator) WriteHeader(statusCode int) {

	// write status code using original http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)

	// capture status code
	r.StatusCode = statusCode
}

func isServerErrorCode(code int) bool {
	return code >= 500 && code <= 599
}
