package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs details of each API request, including response time, headers, and body.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Read request body (for logging)
		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore body after reading

		// Capture response status
		wr := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK, body: &bytes.Buffer{}}

		log.Println("========================================")
		log.Printf("[REQUEST] %s %s", r.Method, r.URL.Path)
		log.Printf("[QUERY PARAMS] %s", r.URL.RawQuery)
		log.Printf("[REMOTE IP] %s", r.RemoteAddr)
		log.Printf("[HEADERS] %v", r.Header)
		if len(body) > 0 {
			log.Printf("[BODY] %s", string(body))
		} else {
			log.Println("[BODY] (empty)")
		}
		log.Println("----------------------------------------")

		next.ServeHTTP(wr, r) // Call the next handler

		duration := time.Since(start)

		log.Println("========================================")
		log.Printf("[RESPONSE] %s %s - %d", r.Method, r.URL.Path, wr.statusCode)
		log.Printf("[RESPONSE TIME] %s", duration)
		log.Printf("[RESPONSE SIZE] %d bytes", wr.body.Len())
		log.Println("========================================")
	})
}

// responseWriter captures the response status and body for logging.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b) // Capture response body
	return rw.ResponseWriter.Write(b)
}
