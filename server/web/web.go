package web

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"server/db"

	"github.com/go-chi/chi/v5"
	"github.com/tdewolff/minify"
)

//go:embed static/*
var staticFS embed.FS

func Router(D *sql.DB, Q db.Queries) chi.Router {
	web := chi.NewRouter()
	web.HandleFunc("/", Home)
	web.HandleFunc("/book", BookHandler(D, Q))
	web.Mount("/static", staticRouter(0))
	return web
}

func minified() http.Handler {
	staticFs, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	m := minify.New()
	sfs := http.FileServer(http.FS(staticFs))
	mini := m.Middleware(http.StripPrefix("/static/", sfs))
	return mini
}
func staticRouter(cacheDuration int) *chi.Mux {
	static := chi.NewRouter()

	// static middleware
	static.Use(
		SetCache(cacheDuration),
	)

	static.Handle("/*", minified())
	return static
}
func SetCache(duration int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", duration))
			next.ServeHTTP(rw, r)
		})
	}
}
