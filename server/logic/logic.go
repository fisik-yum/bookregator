package logic

import (
	"context"
	"database/sql"
	"log"
	"server/db"
)

func ReviewSingle(D *sql.DB, Q db.Queries, ctx context.Context, extrrev db.InsertReviewParams) error{

	err := Q.InsertReview(ctx, extrrev)
	if err != nil {
		return err
	}
	log.Printf("Review for Book %s; User %s inserted", extrrev.Olid, extrrev.Username)
	// update statistics
	s, err := Q.RawStatsFromTable(ctx, extrrev.Olid)
	if err != nil {
		return err
	}
	err = Q.InsertStat(ctx, db.InsertStatParams(computeStats(s)))

	if err != nil {
		return err
	}
	return nil
}

func ReviewMultiple(D *sql.DB, Q db.Queries, ctx context.Context, reviews []db.InsertReviewParams) error{
// Start a transaction
		tx, err := D.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		qtx := Q.WithTx(tx)

		// set
		olidmap := make(map[string]struct{})
		for _, review := range reviews {
			err = qtx.InsertReview(ctx, review)
			if err != nil {
				return err
			}
			olidmap[review.Olid] = struct{}{}
		}

		err = tx.Commit()
		if err != nil {
			return err
		}

		// TODO: reduce code duplication, optimize
		tx, err = D.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		qtx = Q.WithTx(tx)

		// update statistics
		for olid := range olidmap {
			s, err := qtx.RawStatsFromTable(ctx, olid)
			if err != nil {
				return err
			}

			qtx.InsertStat(ctx, db.InsertStatParams(computeStats(s)))
		}
		tx.Commit()
		return nil

}
