/*
TODO: Some form of actual configuration. Use kong? Or roll a custom impl?
*/
package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/api"
	"server/db"
	"server/search"
	"server/web"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var D *sql.DB
var Q db.Queries
var I *search.SearchMachine

// configuration variables

var databaseLocation string // database location
var skipSearchIndex bool // whether to skip indexing
var basicAuthUser string // http basic auth username
var basicAuthPass string // http basic auth password


func init() {
	flag.StringVar(&databaseLocation, "db", "", "location of database file")
	flag.BoolVar(&skipSearchIndex, "skipindex", false, "skip search indexing")
	flag.StringVar(&basicAuthUser, "user", "admin", "http basic auth username")
	flag.StringVar(&basicAuthPass, "pass", "opensesame", "http basic auth password")

	flag.Parse()

	// initialize DB connection
	log.Println("Opening database connection")
	D, Q = db.DBinit(databaseLocation)
	// open up bleve index
	log.Println("Opening bleve index")
	index, err := search.NewSearchMachine("search.bleve")
	if err != nil {
		log.Panic(err)
	}

	I = index
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)

	// manage global mux
	r.Mount("/api", api.Router(D, Q,basicAuthUser, basicAuthPass))
	r.Mount("/", web.Router(D, Q, *I))

	// Bind only to localhost (127.0.0.1)
	addr := "127.0.0.1:1024"
	go func() {
		log.Printf("bookregator listening on http://%s", addr)
		if err := http.ListenAndServe(addr, r); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	if !skipSearchIndex {
		// run refresh in parallel
		ticker := time.NewTicker(30 * time.Minute)
		go func() {
			for {
				log.Println("Refreshing bleve index")
				I.Refresh(D)
				log.Println("Refreshed bleve index")
				<-ticker.C
			}
		}()
		defer ticker.Stop()
	}

	defer I.Close()

	// Ensure that we can shutdown gracefully
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
