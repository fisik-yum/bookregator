package logic

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math"
	"server/db"
	"strings"

	gisbn "github.com/moraes/isbn"
	"gonum.org/v1/gonum/stat/distuv"
)

func computeStats(datum db.RawStatsFromTableRow) db.Stat {
	if datum.CountRatings < 1 {
		return db.Stat{}
	}
	stddev := math.Sqrt(*datum.SumRatingsSquared / float64(datum.CountRatings-1))
	stderror := stddev / math.Sqrt(float64(datum.CountRatings))

	dist := distuv.StudentsT{
		Mu:    *datum.AvgRatings,
		Sigma: stddev,
		Nu:    float64(datum.CountRatings) - 1,
	}

	bound := dist.Quantile(0.975) * stderror
	return db.Stat{
		Olid:        datum.Olid,
		ReviewCount: &datum.CountRatings,
		Rating:      datum.AvgRatings,
		CiBound:     &bound,
	}
}

func InsertRoute(D *sql.DB, Q db.Queries, ctx context.Context, val *db.InsertISBNParams) error {

	// write to DB
	// TODO: Add other sanitizing steps to the scraper layer
	// Remove spaces and dashes
	val.Isbn = strings.ReplaceAll(strings.Trim(val.Isbn, " \n\r"), "-", "")
	val.Olid = strings.ReplaceAll(strings.Trim(val.Olid, " \n\r"), "-", "")

	if !gisbn.Validate(val.Isbn) {
		return errors.New("Invalid ISBN")
	}
	if len(val.Isbn) < 11 && len(val.Isbn) > 8 {
		val.Isbn, _ = gisbn.To13(val.Isbn)
	}
	// get context
	err := Q.InsertISBN(ctx, *val)

	if err != nil {
		return err
	}
	log.Printf("Book routed successfully: %s -> %s", val.Isbn, val.Olid)
	return nil
}

func MassRefreshStats(D *sql.DB, Q db.Queries, ctx context.Context) error {
	rows, err := D.QueryContext(ctx, `SELECT
    		olid AS olid,
    		COUNT(rating) AS count_ratings,
    		AVG(rating) AS avg_ratings,
    		SUM(rating * rating) AS sum_ratings_squared
			FROM reviews GROUP BY olid;
		`)
	if err != nil {
		log.Println(err)
	}

	// run transaction for atomicity and performance
	tx, err := D.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := Q.WithTx(tx)

	for rows.Next() {
		s := db.RawStatsFromTableRow{}
		err = rows.Scan(&s.Olid, &s.CountRatings, &s.AvgRatings, &s.SumRatingsSquared)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(s.Olid)
		err = qtx.InsertStat(ctx, db.InsertStatParams(computeStats(s)))
		if err != nil {
			log.Println(err)
			continue
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func Work(D *sql.DB, Q db.Queries, ctx context.Context, val *db.InsertWorkParams) error {
	err := Q.InsertWork(ctx, *val)
	if err != nil {
		return err
	}
	return nil
}
