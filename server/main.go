package main

import (
	"database/sql"
	"log"
	"net/http"
	"server/handlers/api"
	"server/handlers/web"

	"server/db"
)

var D *sql.DB
var Q db.Queries

func init() {
	D, Q = db.DBinit()
}
func main() {
	globalmux := http.NewServeMux()
	apimux := http.NewServeMux()
	webmux := http.NewServeMux()
	// register apimux endpoints
	apimux.HandleFunc("/api/insert/route", api.InsertRouteHandler(D, Q))
	apimux.HandleFunc("/api/insert/reviewsingle", api.InsertReviewSingleHandler(D, Q))
	apimux.HandleFunc("/api/insert/reviewmultiple", api.InsertReviewMultipleHandler(D, Q))
	apimux.HandleFunc("/api/insert/work", api.InsertWorkHandler(D, Q))
	apimux.HandleFunc("/api/get/reviews", api.GetReviewsHandler(D, Q))
	apimux.HandleFunc("/api/get/work", api.GetWorkHandler(D, Q))
	// register webmux
	webmux.HandleFunc("/book/", web.BookHandler(D, Q))

	// manage global mux
	globalmux.Handle("/api/", apimux)
	globalmux.Handle("/",webmux)
	// Function to log request paths
	handler := Logger{globalmux}

	// Bind only to localhost (127.0.0.1)
	addr := "127.0.0.1:1024"
	log.Printf("bookregator listening on %s", addr)
	if err := http.ListenAndServe(addr, &handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
