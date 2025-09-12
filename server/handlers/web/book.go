package web

import (
	"database/sql"
	"net/http"

	"server/db"
	"server/htmlbuilder/pages"
	"server/htmlbuilder/shared"
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

		work, err := Q.GetWorkByOLID(r.Context(), olid)
		if err != nil {
			return
		}

		pages.NewReview(shared.NewBase(), reviews,work).Render(w, r)
		// render page
	}
}

func Home(w http.ResponseWriter, r *http.Request) {

	pages.NewIndex(shared.NewBase()).Render(w, r)
}
