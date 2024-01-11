// logging.go

package main

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Logging
		end := time.Now()
		duration := end.Sub(start)
		log.Printf("[%s] %s %s %v\n", end.Format("2006-01-02 15:04:05"), r.Method, r.URL.Path, duration)
	})
}
