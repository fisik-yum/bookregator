package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	gisbn "github.com/moraes/isbn"
	"server/db"
)

// single book insert mechanism
func InsertWorkHandler(D *sql.DB, Q db.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		err = Q.InsertWork(ctx, *val)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "DB Write Failed")
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
		// get context
		ctx := r.Context()
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
		err = Q.InsertReview(ctx, extrrev)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "DB Write Failed")
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
		// get context
		ctx := r.Context()
		// Prep data
		reviews := make([]db.InsertReviewParams, 0)
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "JSON Read Failed")
			return
		}
		json.Unmarshal(body, &reviews)

		// Start a transaction
		tx, err := D.Begin()
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "DB Txn start failed")
			return
		}
		defer tx.Rollback()

		qtx := Q.WithTx(tx)
		for _, review := range reviews {
			qtx.InsertReview(ctx, review)
		}
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "DB Txn commit failed")
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

		// write to DB
		// TODO: Add other sanitizing steps to the scraper layer
		// Remove spaces and dashes
		val.Isbn = strings.ReplaceAll(strings.Trim(val.Isbn, " \n\r"), "-", "")
		val.Olid = strings.ReplaceAll(strings.Trim(val.Olid, " \n\r"), "-", "")

		if !gisbn.Validate(val.Isbn) {
			log.Printf("Invalid ISBN: %s", val.Isbn)
			return
		}
		if len(val.Isbn) < 11 && len(val.Isbn) > 8 {
			val.Isbn, _ = gisbn.To13(val.Isbn)
		}
		// get context
		err = Q.InsertISBN(r.Context(), *val)

		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "DB Write Failed")
			return
		}
		log.Printf("Book routed successfully: %s -> %s", val.Isbn, val.Olid)
	}
}
