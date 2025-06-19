package handlers

import (
	"context"
	"database/sql"
	 _"embed"
	"git.sr.ht/~timharek/openlibrary-go"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

var olib openlibrary.Client
var database *sql.DB
var ctx = context.Background()

//go:embed schema.sql
var scheme string

func init() {
	olib = openlibrary.New()

	var err error
	database, err = sql.Open("sqlite3", "test.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	err = database.Ping()
	if err != nil {
		log.Fatal(err)
	}
	// See "Important settings" section.
	database.SetConnMaxLifetime(time.Minute * 3)
	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(10)
	//
	database.ExecContext(ctx, scheme)
}
