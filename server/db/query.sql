-- name: InsertISBN :exec
INSERT INTO isbns (isbn, olid) VALUES (?, ?);

-- name: InsertReview :exec
INSERT INTO reviews ( olid, source, external_id, username, rating, text) values (?, ?, ?, ?, ?, ?);

-- name: InsertWork :exec
INSERT INTO works (olid, title, cover,author, description) values (?, ?, ?, ?, ?);

-- name: GetNReviewsByOLID :many
SELECT * FROM reviews WHERE olid = ? ORDER BY RANDOM() LIMIT ?;

-- name: GetOLIDFromISBN :one
SELECT olid FROM isbns WHERE isbn = ? LIMIT 1;

-- name: GetWorkByOLID :one
SELECT * FROM works WHERE olid = ? LIMIT 1;

-- name: GetStats :one
SELECT * FROM stats WHERE olid = ? LIMIT 1;

-- name: GetRandomWork :one
SELECT olid FROM works ORDER BY RANDOM() LIMIT 1;

-- name: RawStatsFromTable :many
SELECT
    rating
FROM reviews
WHERE
    olid = sqlc.arg(olid) AND rating != -1;

-- name: InsertStat :exec
INSERT  INTO stats (olid, review_count, avg_rating, med_rating, ci_bound) VALUES (?, ?, ?, ?, ?)
ON CONFLICT(olid) DO UPDATE SET review_count=excluded.review_count, avg_rating=excluded.avg_rating, med_rating=excluded.med_rating, ci_bound=excluded.ci_bound
;

-- name: GetISBNRoute :one
SELECT COUNT(*) FROM isbns WHERE isbn = 1 LIMIT 1;
