package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"server/db"
	"server/handlers/api"
	"server/handlers/web"
)

var D *sql.DB
var Q db.Queries

func init() {
	D, Q = db.DBinit()
}
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)

	// manage global mux
	r.Mount("/api", api.Router(D, Q))
	r.Mount("/", web.Router(D,Q))

	// Bind only to localhost (127.0.0.1)
	addr := "127.0.0.1:1024"
	log.Printf("bookregator listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
