package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIRoutesIntegration(t *testing.T) {
	router := setupRouter()

	// Test adding a word
	addWordPayload := `{"Word": "test", "Definition": "definition"}`
	reqAdd := httptest.NewRequest("POST", "/dictionary", bytes.NewBufferString(addWordPayload))
	reqAdd.Header.Set("Content-Type", "application/json")
	rrAdd := httptest.NewRecorder()
	router.ServeHTTP(rrAdd, reqAdd)

	if rrAdd.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rrAdd.Code)
	}

	// Test retrieving a definition
	reqGet := httptest.NewRequest("GET", "/dictionary/test", nil)
	rrGet := httptest.NewRecorder()
	router.ServeHTTP(rrGet, reqGet)

	if rrGet.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rrGet.Code)
	}

	// Test deleting a word
	reqDelete := httptest.NewRequest("DELETE", "/dictionary/test", nil)
	rrDelete := httptest.NewRecorder()
	router.ServeHTTP(rrDelete, reqDelete)

	if rrDelete.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rrDelete.Code)
	}

	// Test listing all entries
	reqList := httptest.NewRequest("GET", "/dictionary", nil)
	rrList := httptest.NewRecorder()
	router.ServeHTTP(rrList, reqList)

	if rrList.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rrList.Code)
	}

	// Test authentication middleware
	reqAuth := httptest.NewRequest("GET", "/dictionary", nil)
	reqAuth.Header.Set("Authorization", "Bearer XdkmKQ79Q979XxEVHcUS1liWoDUhPACF1DuAGWUqPSxlWRr3JHr86l2okqOh")
	rrAuth := httptest.NewRecorder()
	router.ServeHTTP(rrAuth, reqAuth)

	if rrAuth.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rrAuth.Code)
	}

	// Add more test cases as needed
}
