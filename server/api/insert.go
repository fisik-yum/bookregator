package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"server/db"
	"server/logic"
)

// single book insert mechanism
func InsertWorkHandler(D *sql.DB, Q db.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "JSON Read Failed")
			return
		}
		val := new(db.InsertWorkParams)
		json.Unmarshal(raw, val)

		// write to DB
		err = logic.Work(D,Q,r.Context(),val)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// Single review insert mechanism
func InsertReviewSingleHandler(D *sql.DB, Q db.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// prep data
		review := new(db.Review)
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "JSON Read Failed")
			return
		}
		json.Unmarshal(body, review)
		extrrev := db.InsertReviewParams{
			Olid:       review.Olid,
			Source:     review.Source,
			ExternalID: review.ExternalID,
			Username:   review.Username,
			// fix these two
			Rating: (review.Rating),
			Text:   review.Text,
		}

		err = logic.ReviewSingle(D, Q, r.Context(), extrrev)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Review for Book %s; User %s inserted", review.Olid, review.Username)
	}
}

// Multiple review insert mechanism
func InsertReviewMultipleHandler(D *sql.DB, Q db.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// transactions with sqlc
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Prep data
		reviews := make([]db.InsertReviewParams, 0)
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		json.Unmarshal(body, &reviews)

		err = logic.ReviewMultiple(D, Q, r.Context(), reviews)
		if err != nil {
			log.Println(err)
			return
		}

	}
}

/*
this handler simply exists for debugging. we force a refresh of all statistics
*/
func UpdateStatGlobal(D *sql.DB, Q db.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := logic.MassRefreshStats(D, Q, r.Context())
		if err != nil {
			log.Println(err)
			return
		}
	}
}

/*
Extract ISBN, and auto-routes it. Data is sent as JSON,
*/
func InsertRouteHandler(D *sql.DB, Q db.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "JSON Read Failed")
			return
		}
		val := new(db.InsertISBNParams)
		json.Unmarshal(raw, val)

		err = logic.InsertRoute(D, Q, r.Context(), val)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Book routed successfully: %s -> %s", val.Isbn, val.Olid)
	}
}
