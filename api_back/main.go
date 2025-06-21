package main

import (
	"api_back/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/insert/route", handlers.InsertRouteHandler)
	http.HandleFunc("/api/insert/reviewsingle", handlers.InsertReviewSingleHandler)
	http.HandleFunc("/api/insert/reviewmultiple", handlers.InsertReviewMultipleHandler)
	// Bind only to localhost (127.0.0.1)
	addr := "127.0.0.1:8080"
	log.Printf("Ingestion API listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}


