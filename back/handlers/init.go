package handlers

import (
	"database/sql"

	"book.buckminsterfullerene.net/db"
)

var D *sql.DB
var Q db.Queries

func init() {
	D, Q = db.DBinit()
}
