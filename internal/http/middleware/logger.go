package middleware

import (
	"log"
	"net/http"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(recorder, r)

		// Log the request
		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.URL.Path,
			http.StatusText(recorder.Status),
			time.Since(start),
		)
	})
}
