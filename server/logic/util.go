package logic

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"server/db"
	"sort"
	"strings"

	gisbn "github.com/moraes/isbn"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

func computeStats(datum []float64) db.Stat {
	if len(datum) < 1 {
		return db.Stat{
			ReviewCount: 0,
			AvgRating: -1,
			MedRating: -1,
			CiBound: -1,
		}
	}

	sort.Float64s(datum)

	review_count := float64(len(datum))
	mean, stddev := stat.MeanStdDev(datum, nil)
	stderror := stat.StdErr(stddev, review_count)
	median := stat.Quantile(0.5, stat.Empirical, datum, nil)

	dist := distuv.StudentsT{
		Mu:    mean,
		Sigma: stderror,
		Nu:    review_count - 1,
	}

	bound := dist.Quantile(0.975) * stderror
	return db.Stat{
		ReviewCount:           int64(review_count),
		AvgRating:             mean,
		MedRating:             median,
		CiBound:               bound,
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
    		olid FROM works`)
	if err != nil {
		log.Println(err)
	}

	var olids []string
	for rows.Next() {
		var olid string
		if err := rows.Scan(&olid); err != nil {
			return err
		}
		olids = append(olids, olid)
	}
	if err := rows.Close(); err != nil {
		return err
	}

	// run transaction for atomicity and performance
	tx, err := D.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := Q.WithTx(tx)

	for _, k := range olids {
		d, err := qtx.RawStatsFromTable(ctx, k)
		if err != nil {
			return err
		}
		s := computeStats(d)
		s.Olid = k
		qtx.InsertStat(ctx, db.InsertStatParams(s))
		log.Printf("Updated statistics for book %s",k)
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
