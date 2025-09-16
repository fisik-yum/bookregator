package api

import (
	"database/sql"
	"server/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router(D *sql.DB, Q db.Queries) chi.Router {
	api := chi.NewRouter()
	api.Mount("/insert", secureRouter(D, Q))
	api.Route("/get", func(r chi.Router) {
		r.Get("/reviews", GetReviewsHandler(D, Q))
		r.Get("/work", GetWorkHandler(D, Q))
	})
	return api
}

func secureRouter(D *sql.DB, Q db.Queries) chi.Router {
	// TODO: fix "secure" -> actually add config
	secure := chi.NewRouter().With(middleware.BasicAuth("Insert Endpoints", map[string]string{"scraper": "opensesame"}))

	secure.Route("/", func(r chi.Router) {
		r.Post("/route", InsertRouteHandler(D, Q))
		r.Post("/work", InsertWorkHandler(D, Q))
		r.Post("/reviewsingle", InsertReviewSingleHandler(D, Q))
		r.Post("/reviewmultiple", InsertReviewMultipleHandler(D, Q))
	})

	return secure
}
