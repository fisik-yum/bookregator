/*
TODO: Some form of actual configuration. Use kong? Or roll a custom impl?
*/
package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/api"
	"server/db"
	"server/search"
	"server/web"
	"time"

	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var D *sql.DB
var Q db.Queries
var I *search.SearchMachine

func init() {
	// initialize DB connection
	log.Println("Opening database connection")
	D, Q = db.DBinit("test.sqlite3")
	// open up bleve index
	log.Println("Opening bleve index")
	index, err := search.NewSearchMachine("test.bleve")
	if err != nil {
		log.Panic(err)
	}

	I = index
	// index all items on startup
	log.Println("Refreshing bleve index")
	I.Refresh(D)
	log.Println("Refreshed bleve index")

	log.Println("Running test query")
	results, err := I.SearchItem("ancillary", context.Background())
	if err != nil {
		log.Panic(err)
	}
	for _, v := range results {
		log.Printf("%s", v)
	}

}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)

	// manage global mux
	r.Mount("/api", api.Router(D, Q))
	r.Mount("/", web.Router(D, Q, *I))

	// Bind only to localhost (127.0.0.1)
	addr := "127.0.0.1:1024"
	go func() {
		log.Printf("bookregator listening on %s", addr)
		if err := http.ListenAndServe(addr, r); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// run refresh in parallel
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for {
			<-ticker.C
			log.Println("Refreshing bleve index")
			I.Refresh(D)
			log.Println("Refreshed bleve index")
		}
	}()

	defer I.Close()
	defer ticker.Stop()

	// Ensure that we can shutdown gracefully
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
