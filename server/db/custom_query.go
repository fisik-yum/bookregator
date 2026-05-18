package db

import (
	"context"
)

func (q *Queries) PruneDatabase(ctx context.Context) error {
	// prune works
	var err error
	_, err = q.db.ExecContext(ctx, `DELETE FROM works
    WHERE olid in (
        select w2.olid from works w2 left join reviews r
        on r.olid=w2.olid group by w2.olid 
        having count(r.source) = 0
    );`,
	)
	if err != nil {
		return err
	}
	// prune reviews
	_, err = q.db.ExecContext(ctx, `DELETE FROM reviews where 
    WHERE olid not in (select w.olid from works w);`,
	)
	if err != nil {
		return err
	}
	// prune stats
	_, err = q.db.ExecContext(ctx, `DELETE FROM stats where 
    WHERE olid not in (select w.olid from works w);`,
	)
	if err != nil {
		return err
	}
	// prune routing
	_, err = q.db.ExecContext(ctx, `DELETE FROM isbns where 
    WHERE olid not in (select w.olid from works w);`,
	)
	if err != nil {
		return err
	}
	return nil
}
