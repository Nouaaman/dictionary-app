package middleware

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	// Open or create a log file
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}

	logger := log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logString := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		// Write logs to both console and file
		logger.Println(logString)
		defer logFile.Close()
	})

}
