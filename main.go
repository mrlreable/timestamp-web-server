package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Make a buffered channel (with size 1) for handling concurrent read/write operations
var c = make(chan time.Time, 1)

func fetchTime(w http.ResponseWriter, req *http.Request) {
	fmt.Println("GET /time/fetch handler called")

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	// Read the value from the channel and put it back again
	ts := <-c
	c <- ts

	response := fmt.Sprintf("Time: %s", ts.Format(time.UnixDate))
	w.Write([]byte(response))
}

func setTime(w http.ResponseWriter, req *http.Request) {
	fmt.Println("POST /time/set handler called")

	// Read the body and return if provided value is invalid
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Malformed request", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	parsedTime, err := time.Parse(time.UnixDate, string(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid Unix time %s", body), http.StatusBadRequest)
		return
	}

	// Remove the time from the channel and push the parsed time
	<-c
	c <- parsedTime

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Stored timestamp"))
}

func main() {
	fmt.Println("Timestamp server spinning up...")

	http.Handle("/time/set", middleware(http.MethodPost, http.HandlerFunc(setTime)))
	http.Handle("/time/fetch", middleware(http.MethodGet, http.HandlerFunc(fetchTime)))

	// Init channel with an empty timestamp
	c <- time.Time{}

	// Simulate client
	go simulateClient()

	// Start server on a goroutine
	fmt.Println("Listening on http://127.0.0.1:8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		fmt.Println("Error starting server on http://127.0.0.1:8090 Error:", err)
	}
}

// Utils
// Ideally separated into different package
func middleware(m string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Content-Type") != "text/plain" && req.Method == http.MethodPost {
			fmt.Println("Content-Type:", req.Header.Get("Content-Type"))
			http.Error(w, "Unsupported Media Type. Expected 'text/plain'", http.StatusUnsupportedMediaType)
			return
		}
		if req.Method != m {
			http.Error(w, fmt.Sprintf("Method not allowed. Expected '%s'", m), http.StatusMethodNotAllowed)
		}
		next.ServeHTTP(w, req)
	})
}

// Simulate client
func simulateClient() {
	// Wait for server to start up
	// Not ideal but will work in this scenario
	time.Sleep(2 * time.Second)

	timestamp := time.Now().Format(time.UnixDate)

	// Test POST /time/set
	resp, err := http.Post("http://127.0.0.1:8090/time/set", "text/plain", strings.NewReader(timestamp))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	// Test GET /time/fetch
	resp, err = http.Get("http://127.0.0.1:8090/time/fetch")
	if err != nil {
		fmt.Println("Error sending GET request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println(string(body))
}
