package api

import (
	"database/sql"
	"server/db"

	"github.com/go-chi/chi/v5"
)

func Router(D *sql.DB, Q db.Queries) chi.Router {
	api:=chi.NewRouter()
	api.Route("/insert", func(r chi.Router){
		r.Post("/route",InsertRouteHandler(D, Q))
		r.Post("/work",InsertWorkHandler(D, Q))
		r.Post("/reviewsingle",InsertReviewSingleHandler(D, Q))
		r.Post("/reviewmultiple",InsertReviewMultipleHandler(D, Q))
	})
	api.Route("/get", func(r chi.Router){
		r.Get("/reviews",GetReviewsHandler(D, Q))
		r.Get("/work",GetWorkHandler(D, Q))
	})
	return api
}
