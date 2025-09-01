package handlers

// TODO: Better logging
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


// Single review insert mechanism
func InsertReviewSingleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// get context
	ctx := r.Context()
	// prep data
	review := new(common.Review)
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
		Rating:     &(review.Rating),
		Text:       &review.Text,
	}
	err = internal.Queries.InsertReview(ctx, extrrev)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "DB Write Failed")
	}
	log.Printf("Review for Book %s; User %s inserted", review.Olid, review.Username)
}

// Multiple review insert mechanism
func InsertReviewMultipleHandler(w http.ResponseWriter, r *http.Request) {
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
	tx, err := internal.Database.Begin()
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "DB Txn start failed")
		return
	}
	defer tx.Rollback()

	qtx := internal.Queries.WithTx(tx)
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
