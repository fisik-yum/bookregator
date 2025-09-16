package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"

	"server/db"

	gisbn "github.com/moraes/isbn"
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
			return
		}
		log.Printf("Review for Book %s; User %s inserted", review.Olid, review.Username)
		// update statistics
		s, err := Q.RawStatsFromTable(ctx, extrrev.Olid)
		if err != nil {
			log.Println(err)
			return
		}
		err = Q.InsertStat(ctx, db.InsertStatParams(statsCalculator(D,Q,s)))

		// commit
		if err != nil {
			log.Println(err)
			return
		}
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
			return
		}
		json.Unmarshal(body, &reviews)

		// Start a transaction
		tx, err := D.Begin()
		if err != nil {
			log.Println(err)
			return
		}
		defer tx.Rollback()

		qtx := Q.WithTx(tx)

		// set
		olidmap := make(map[string]struct{})
		for _, review := range reviews {
			err = qtx.InsertReview(ctx, review)
			if err != nil {
				log.Println(err)
				return
			}
			olidmap[review.Olid] = struct{}{}
		}

		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return
		}

		// TODO: reduce code duplication, optimize
		tx, err = D.Begin()
		if err != nil {
			log.Println(err)
			return
		}
		defer tx.Rollback()

		qtx = Q.WithTx(tx)

		// update statistics
		for olid := range olidmap {
			s, err := qtx.RawStatsFromTable(ctx, olid)
			if err != nil {
				log.Println(err)
				return
			}

			qtx.InsertStat(ctx, db.InsertStatParams(statsCalculator(D,Q,s)))
		}
		tx.Commit()
	}
}

/*
Extract ISBN, and auto-routes it. Data is sent as JSON,
*/func InsertRouteHandler(D *sql.DB, Q db.Queries) func(w http.ResponseWriter, r *http.Request) {
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

func statsCalculator(D *sql.DB, Q db.Queries, datum db.RawStatsFromTableRow) db.Stat {
	if datum.CountRatings < 1 {
		return db.Stat{}
	}
	stddev := math.Sqrt(*datum.SumRatingsSquared / float64(datum.CountRatings-1))
	stderror := stddev / math.Sqrt(float64(datum.CountRatings))
	bound := 1.96 * stderror
	return db.Stat{
		Olid:        datum.Olid,
		ReviewCount: &datum.CountRatings,
		Rating:      datum.AvgRatings,
		CiBound:     &bound,
	}
}
