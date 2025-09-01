package handlers

import (
	"api_back/internal"
	"api_back/internal/db"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "git.sr.ht/~timharek/openlibrary-go"
	gisbn "github.com/moraes/isbn"
)

/*
Extract ISBN, and auto-routes it. Requires URL form to have the fields "isbn"
and "olid". ISBNs can be ISBN10 or ISBN13.
Example: endpoint?isbn=123456790?olid=123456
*/
func InsertRouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// prep data
	r.ParseForm()
	isbn := r.Form.Get("isbn")
	if isbn == "" {
		fmt.Fprintln(w, "No ISBN in request URL")
		return
	}
	olid := r.Form.Get("olid")
	if olid == "" {
		fmt.Fprintln(w, "No OLID in request URL")
		return
	}
	if routeISBNtoOLID(isbn, olid, r) {
		fmt.Fprintf(w, "Success: Route %s to %s\n", isbn, olid)
		log.Printf("Success: Route %s to %s", isbn, olid)
	} else {
		fmt.Fprintf(w, "Fail: Route %s to %s\n", isbn, olid)
		log.Printf("Fail: Route %s to %s", isbn, olid)
	}
}

/*
Patches ISBN to an OLID and adds it to the database. Performs only input
validation, not verifications. That is is done on the client side.
*/
func routeISBNtoOLID(isbn string, olid string, r *http.Request) bool {
	// TODO: Add other sanitizing steps to the scraper layer
	// Remove spaces and dashes
	isbn = strings.ReplaceAll(strings.Trim(isbn, " \n\r"), "-", "")
	olid = strings.ReplaceAll(strings.Trim(olid, " \n\r"), "-", "")

	if !gisbn.Validate(isbn) {
		log.Printf("Invalid ISBN: %s", isbn)
		return false
	}
	if len(isbn) < 11 && len(isbn) > 8 {
		isbn, _ = gisbn.To13(isbn)
	}
	// get context
	ctx := r.Context()
	err := internal.Queries.InsertISBN(ctx, db.InsertISBNParams{
		Isbn: isbn,
		Olid: olid,
	})

	if err != nil {
		return false
	}
	return true
}
