package web

import (
	"database/sql"
	"server/db"

	"github.com/go-chi/chi/v5"
)

func Router(D *sql.DB, Q db.Queries) chi.Router {
	web := chi.NewRouter()
	web.HandleFunc("/", Home)
	web.HandleFunc("/book", BookHandler(D, Q))
	return web
}
