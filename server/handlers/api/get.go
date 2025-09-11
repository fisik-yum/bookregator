package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"book.buckminsterfullerene.net/db"
)

// get book reviews
func GetReviewsHandler(D *sql.DB, Q db.Queries) func(w http.ResponseWriter, r *http.Request) {
	// TODO: Move this to frontend
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "JSON Read Failed")
			return
		}
		val := new(db.GetXByOLIDParams)
		json.Unmarshal(raw, val)

		reviews, err := Q.GetNReviewsByOLID(r.Context(), db.GetNReviewsByOLIDParams{
			Olid:  val.OLID,
			Limit: 5,
		})
		if err != nil {
			log.Printf("Review Request for OLID: %s failed", val.OLID)
			log.Println(err)
		}

		reviewCommon := make([]db.Review, len(reviews))
		for i, v := range reviews {
			rev := db.Review{
				Olid:       v.Olid,
				Text:       v.Text,
				Source:     v.Source,
				Rating:     v.Rating,
				Username:   v.Username,
				ExternalID: v.Username,
			}
			reviewCommon[i] = rev
		}

		reviewJSON, err := json.Marshal(reviewCommon)
		if err != nil {
			log.Printf("Review Request for OLID: %s failed", val.OLID)
			log.Println(err)
		}
		w.Write(reviewJSON)
	}
}

func GetWorkHandler(D *sql.DB, Q db.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "JSON Read Failed")
			return
		}

		val := new(db.GetXByOLIDParams)
		json.Unmarshal(raw, val)

		work, err := Q.GetWorkByOLID(r.Context(), val.OLID)

		workJSON, err := json.Marshal(work)
		if err != nil {
			log.Printf("Work Request for OLID: %s failed", val.OLID)
			log.Println(err)
		}
		w.Write(workJSON)

	}
}
