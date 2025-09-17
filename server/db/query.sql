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

-- name: RawStatsFromTable :one
SELECT
    olid AS olid,
    COUNT(rating) AS count_ratings,
    AVG(rating) AS avg_ratings,
    SUM(rating * rating) AS sum_ratings_squared
FROM reviews
WHERE
    olid = sqlc.arg(olid) AND rating != -1;

-- name: InsertStat :exec
INSERT  INTO stats (olid, review_count, rating, ci_bound) VALUES (?, ?, ?, ?)
ON CONFLICT(olid) DO UPDATE SET review_count=excluded.review_count, rating=excluded.rating, ci_bound=excluded.ci_bound
;
