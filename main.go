package main

import (
	"fmt"
	"net/http"
)

func fetchTime(w http.ResponseWriter, req *http.Request) {
	fmt.Println("GET /time/fetch handler called")
}

func setTime(w http.ResponseWriter, req *http.Request) {
	fmt.Println("POST /time/set handler called")
}

func main() {
	fmt.Println("Timestamp server spinning up...")

	http.Handle("/time/set", middleware(http.MethodPost, http.HandlerFunc(setTime)))
	http.Handle("/time/fetch", middleware(http.MethodGet, http.HandlerFunc(fetchTime)))

	fmt.Println("Listening on :8090")
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("Error starting server on :8090 Error:", err)
	}
}

// Utils
// Ideally separated into different package
func middleware(m string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "Unsupported Media Type. Expected 'text/plain'", http.StatusUnsupportedMediaType)
			return
		}
		if req.Method != m {
			http.Error(w, fmt.Sprintf("Method not allowed. Expected '%s'", m), http.StatusMethodNotAllowed)
		}
		next.ServeHTTP(w, req)
	})
}
