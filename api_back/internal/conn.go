package internal

import (
	"api_back/internal/db"
	"context"
	"database/sql"
	_ "embed"
	"log"
	"time"

	"git.sr.ht/~timharek/openlibrary-go"
	_ "github.com/mattn/go-sqlite3"
)

var olib openlibrary.Client
var Database *sql.DB
var Queries db.Queries

//go:embed schema.sql
var scheme string

func init() {
	olib = openlibrary.New()

	var err error
	Database, err = sql.Open("sqlite3", "test.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	err = Database.Ping()
	if err != nil {
		log.Fatal(err)
	}
	// See "Important settings" section.
	Database.SetConnMaxLifetime(time.Minute * 3)
	Database.SetMaxOpenConns(10)
	Database.SetMaxIdleConns(10)
	// initialize
	Database.ExecContext(context.Background(), scheme)
	Queries = *db.New(Database)

}
