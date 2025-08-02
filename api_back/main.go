package main

import (
	"api_back/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	// register
	mux.HandleFunc("/api/insert/route", handlers.InsertRouteHandler)
	mux.HandleFunc("/api/insert/reviewsingle", handlers.InsertReviewSingleHandler)
	mux.HandleFunc("/api/insert/reviewmultiple", handlers.InsertReviewMultipleHandler)
	mux.HandleFunc("/api/insert/work", handlers.InsertWorkHandler)
	mux.HandleFunc("/api/get/work", handlers.GetReviewsHandler)

	// Function to log request paths
	loggingFn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s", r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	}

	// Bind only to localhost (127.0.0.1)
	addr := "127.0.0.1:1024"
	log.Printf("Ingestion API listening on %s", addr)
	if err := http.ListenAndServe(addr, http.HandlerFunc(loggingFn)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
