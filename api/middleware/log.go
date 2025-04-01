package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs details of each API request, including response time.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("API request")
		// Capture response status
		wr := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wr, r)

		duration := time.Since(start)
		log.Printf("[%s] %s %s %d %s", r.Method, r.URL.Path, r.RemoteAddr, wr.statusCode, duration)
	})
}

// responseWriter is a wrapper to capture HTTP status codes.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
