package handlers

import (
	"api_back/db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	err = queries.InsertWork(ctx, *val)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "DB Write Failed")
		return
	}
}
