package handlers
// TODO: Better logging
import (
	"api_back/db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	review := new(db.InsertReviewParams)
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "JSON Read Failed")
		return
	}
	json.Unmarshal(body, review)
	err = queries.InsertReview(ctx, *review)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "DB Write Failed")
	}
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
	tx, err := database.Begin()
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "DB Txn start failed")
		return
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)
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
