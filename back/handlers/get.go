package handlers

import (
	"log"
	"encoding/json"
	"book.buckminsterfullerene.net/db"
	"net/http"
)

// get book reviews
func GetReviewsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	olid := r.Form.Get("olid")
	reviews, err := Q.GetNReviewsByOLID(r.Context(), db.GetNReviewsByOLIDParams{
		Olid:  olid,
		Limit: 7,
	})
	if err != nil {
		log.Printf("Review Request for OLID: %s failed", olid)
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
		log.Printf("Review Request for OLID: %s failed", olid)
		log.Println(err)
	}
	w.Write(reviewJSON)
}
