package web

import (
	"database/sql"
	"net/http"

	"server/db"
	"server/web/pages"
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

		ctx:=r.Context()
		
		work, err := Q.GetWorkByOLID(ctx, olid)
		if err != nil {
			return
		}
		stat,err:=Q.GetStats(ctx,olid)
		if err != nil {
			return
		}

		// render page
		pages.NewReview(reviews,work,stat).Render(w, r)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {

	pages.NewIndex().Render(w, r)
}
