package web

import (
	"database/sql"
	"net/http"

	"server/db"
	"server/htmlbuilder"

)

func BookHandler(D *sql.DB, Q db.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		olid := r.Form.Get("olid")

		// get data to hydrate the page
		reviews, err := Q.GetNReviewsByOLID(r.Context(), db.GetNReviewsByOLIDParams{
			Olid:  olid,
			Limit: 5,
		})
		if err != nil {
			return
		}
		// render page
		v := htmlbuilder.ReviewPage(reviews)
		v.Render(w)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	htmlbuilder.Home().Render(w)
}
