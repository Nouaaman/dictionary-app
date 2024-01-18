package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test with a valid authentication token
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer XdkmKQ79Q979XxEVHcUS1liWoDUhPACF1DuAGWUqPSxlWRr3JHr86l2okqOh")
	rr := httptest.NewRecorder()

	AuthenticationMiddleware(handler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Test with an invalid or missing authentication token
	req = httptest.NewRequest("GET", "/", nil)
	rr = httptest.NewRecorder()

	AuthenticationMiddleware(handler).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test if logging middleware logs the expected information
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	LoggingMiddleware(handler).ServeHTTP(rr, req)

	// Add more test cases as needed
}
