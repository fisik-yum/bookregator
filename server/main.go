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
	"github.com/ilyakaznacheev/cleanenv"
)

var D *sql.DB
var Q db.Queries
var I *search.SearchMachine

// configuration 
var configLocation string // config location
var cfg serverConfig

func init() {
	flag.StringVar(&configLocation, "c", "", "location of config file")
	flag.Parse()

	err:=cleanenv.ReadConfig(configLocation,&cfg)
	if err != nil {
		log.Panic(err)
	}
	log.Println(cfg)

	// initialize DB connection
	log.Println("Opening database connection")
	D, Q = db.DBinit(cfg.DatabaseLocation)
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
	r.Mount("/api", api.Router(D, Q,cfg.HTTPUser, cfg.HTTPPass))
	r.Mount("/", web.Router(D, Q, *I))

	// Bind only to localhost (127.0.0.1)
	addr := "127.0.0.1:1024"
	go func() {
		log.Printf("bookregator listening on http://%s", addr)
		if err := http.ListenAndServe(addr, r); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	if !cfg.SkipIndexing{
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
