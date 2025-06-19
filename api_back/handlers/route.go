package handlers

import (
	"api_back/db"
	"fmt"
	"log"
	"net/http"

	_ "git.sr.ht/~timharek/openlibrary-go"
	gisbn "github.com/moraes/isbn"
)

func RouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	isbn := r.Form.Get("isbn")
	if isbn == "" {
		log.Println("No ISBN in request URL")
		return
	}
	if routeISBNtoOLID(isbn) {
		fmt.Fprintln(w, "Review ingested successfully")
	} else {
		fmt.Fprintln(w, "Review ingested unsuccessfully")
	}
}
func routeISBNtoOLID(isbn string) bool {
	if !gisbn.Validate(isbn) {
		log.Printf("Invalid ISBN: %s", isbn)
		return false
	}
	if len(isbn) < 11 && len(isbn) > 8 {
		isbn, _ = gisbn.To13(isbn)
	}
	//
	book, err := olib.Book.ByISBN(isbn)
	if err != nil {
		log.Printf("Book %s lookup error: %s", isbn, err)
		return false
	}
	olid := book.Key // HACK: Assume that this is the correct OLID for the works
	if olid == "" {
		log.Printf("ISBN %s has no corresponding OLID", isbn)
		return false
	}

	/*_, err = db.Exec(`
	      INSERT OR IGNORE INTO isbns (isbn, olid)
	      VALUES (?, ?)
	  `, isbn, olid)
	*/
	queries := db.New(database)
	err = queries.InsertISBN(ctx, db.InsertISBNParams{
		Isbn: isbn,
		Olid: olid,
	})

	if err != nil {
		log.Println("DB insert failed:", err)
		return false
	}

	log.Printf("DB insert suceeded: %s", olid)
	return true
}
