package web

import (
	"database/sql"
	"log"
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
			log.Println(err)
			return
		}

		ctx:=r.Context()
		
		work, err := Q.GetWorkByOLID(ctx, olid)
		if err != nil {
			log.Println(err)
			return
		}
		stat,err:=Q.GetStats(ctx,olid)
		if err != nil {
			log.Println(err)
			//return
		}

		// render page
		pages.NewReview(reviews,work,stat).Render(w, r)
	}
}


func RandomBookHandler(D *sql.DB, Q db.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r*http.Request){
		ctx:=r.Context()
		olid,err:=Q.GetRandomWork(ctx)
		if err!=nil{
			log.Println(err)
			return
		}
		http.Redirect(w,r,"book?olid="+olid,http.StatusPermanentRedirect)
	}
}
func Home(w http.ResponseWriter, r *http.Request) {

	pages.NewIndex().Render(w, r)
}
