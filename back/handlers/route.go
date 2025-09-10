package handlers

import (
	"api_back/internal"
	"api_back/internal/db"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"encoding/json"
	"io"

	_ "git.sr.ht/~timharek/openlibrary-go"
	gisbn "github.com/moraes/isbn"
)

/*
Extract ISBN, and auto-routes it. Data is sent as JSON,
*/
func InsertRouteHandler(w http.ResponseWriter, r *http.Request) {
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

	if !routeISBNtoOLID(val, r.Context()) {
		log.Println(err)
		fmt.Fprintf(w, "DB Write Failed")
		return
	}
	log.Printf("Book routed successfully: %s -> %s", val.Isbn, val.Olid)
}

/*
Patches ISBN to an OLID and adds it to the database. Performs only input
validation, not verifications. That is is done on the client side.
*/
func routeISBNtoOLID(ins *db.InsertISBNParams, ctx context.Context) bool {
	// TODO: Add other sanitizing steps to the scraper layer
	// Remove spaces and dashes
	ins.Isbn = strings.ReplaceAll(strings.Trim(ins.Isbn, " \n\r"), "-", "")
	ins.Olid = strings.ReplaceAll(strings.Trim(ins.Olid, " \n\r"), "-", "")

	if !gisbn.Validate(ins.Isbn) {
		log.Printf("Invalid ISBN: %s", ins.Isbn)
		return false
	}
	if len(ins.Isbn) < 11 && len(ins.Isbn) > 8 {
		ins.Isbn, _ = gisbn.To13(ins.Isbn)
	}
	// get context
	err := internal.Queries.InsertISBN(ctx, *ins)

	if err != nil {
		return false
	}
	return true
}
