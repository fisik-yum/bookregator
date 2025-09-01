package main

import (
	"log"
	"net/http"
	"api_front/handlers"
)

func main() {
	mux := http.NewServeMux()
	// register
	mux.HandleFunc("/bookregator/olid", handlers.ByOLID)

	// Function to log request paths
	loggingFn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s", r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	}

	// Bind only to localhost (127.0.0.1)
	addr := "127.0.0.1:8080"
	log.Printf("Webservice listening on %s", addr)
	if err := http.ListenAndServe(addr, http.HandlerFunc(loggingFn)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
