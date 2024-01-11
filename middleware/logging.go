// middleware/logging.go

package middleware

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	// Open or create a log file
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()

	// Use multi-writer to write logs to both console and file
	logWriter := io.MultiWriter(os.Stdout, logFile)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Logging
		time := time.Now()
		logString := fmt.Sprintf("[%s] %s %s \n", time.Format("2006-01-02 15:04:05"), r.Method, r.URL.Path)

		// Write logs to both console and file
		log.Println(logString)
		fmt.Fprintln(logWriter, logString)
	})
}
