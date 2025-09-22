package db

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"time"

	"git.sr.ht/~timharek/openlibrary-go"
	_ "github.com/mattn/go-sqlite3"
)

var olib openlibrary.Client

//go:embed schema.sql
var scheme string

func DBinit(loc string) (D *sql.DB, Q Queries) {
	olib = openlibrary.New()

	var err error
	D, err = sql.Open("sqlite3", loc)
	if err != nil {
		log.Fatal(err)
	}
	err = D.Ping()
	if err != nil {
		log.Fatal(err)
	}
	// See "Important settings" section.
	D.SetConnMaxLifetime(time.Minute * 3)
	D.SetMaxOpenConns(10)
	D.SetMaxIdleConns(10)
	// initialize
	// TODO: make this optional, if we want open readonly
	D.ExecContext(context.Background(), scheme)
	Q = *New(D)
	return
}
