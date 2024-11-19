package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSetTimeValidFormat(t *testing.T) {
	// Init channel
	c <- time.Time{}

	// POST /time/set
	timestamp := time.Now().Format(time.UnixDate)
	req := httptest.NewRequest(http.MethodPost, "/time/set", strings.NewReader(timestamp))
	req.Header.Set("Content-Type", "text/plain")

	rr := httptest.NewRecorder()

	setTime(rr, req)

	// Assert response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check if the channel contains the correct value
	receivedTime := <-c
	c <- receivedTime
	parsedTime, _ := time.Parse(time.UnixDate, timestamp)
	if !receivedTime.Equal(parsedTime) {
		t.Errorf("Expected timestamp %v, got %v", parsedTime, receivedTime)
	}
}

func TestSetTimeRandomStringGiven(t *testing.T) {
	// Init channel
	c <- time.Time{}

	// POST /time/set
	invalidTimestamp := "this is an invalid time string"
	req := httptest.NewRequest(http.MethodPost, "/time/set", strings.NewReader(invalidTimestamp))
	req.Header.Set("Content-Type", "text/plain")

	rr := httptest.NewRecorder()

	setTime(rr, req)

	// Assert response status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestSetTimeInvalidTimeFormat(t *testing.T) {
	// Init channel
	c <- time.Time{}

	// Simulate an invalid POST /time/set
	invalidTime := time.Now().Format(time.RFC1123)
	req := httptest.NewRequest(http.MethodPost, "/time/set", strings.NewReader(invalidTime))
	req.Header.Set("Content-Type", "text/plain")

	rr := httptest.NewRecorder()

	setTime(rr, req)

	// Assert response status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestFetchTime(t *testing.T) {
	// Initialize the channel with a specific timestamp
	expectedTime := time.Now()
	c <- expectedTime

	// GET /time/fetch
	req := httptest.NewRequest(http.MethodGet, "/time/fetch", nil)

	rr := httptest.NewRecorder()

	fetchTime(rr, req)

	// Assert response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Assert response body
	expectedResponse := "Time: " + expectedTime.Format(time.UnixDate)
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected body %q, got %q", expectedResponse, rr.Body.String())
	}
}

func TestMiddleware(t *testing.T) {
	// Dummy handler to wrap with middleware
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		method         string
		contentType    string
		expectedStatus int
	}{
		{"Valid POST", http.MethodPost, "text/plain", http.StatusOK},
		{"Invalid Content-Type", http.MethodPost, "application/json", http.StatusUnsupportedMediaType},
		{"Invalid Method", http.MethodGet, "text/plain", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			middleware(http.MethodPost, dummyHandler).ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
