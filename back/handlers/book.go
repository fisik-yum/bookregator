package handlers

import (
	"api_back/internal"
	"api_back/internal/db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"book.buckminsterfullerene.net/common"
)

// single book insert mechanism
func InsertWorkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "JSON Read Failed")
		return
	}
	val := new(db.InsertWorkParams)
	json.Unmarshal(raw, val)

	// write to DB
	err = internal.Queries.InsertWork(ctx, *val)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "DB Write Failed")
		return
	}
}

// get book reviews
func GetReviewsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	olid := r.Form.Get("olid")
	reviews, err := internal.Queries.GetNReviewsByOLID(r.Context(), db.GetNReviewsByOLIDParams{
		Olid:  olid,
		Limit: 7,
	})
	if err != nil {
		log.Printf("Review Request for OLID: %s failed", olid)
		log.Println(err)
	}

	reviewCommon := make([]common.Review, len(reviews))
	for i, v := range reviews {
		rev := common.Review{
			Olid:       v.Olid,
			Text:       *v.Text,
			Source:     v.Source,
			Rating:     *v.Rating,
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
