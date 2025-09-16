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

-- name: UpdateStatistics :exec
INSERT INTO stats (olid, rating, review_count) 
VALUES (sqlc.arg(olid), (SELECT COALESCE(AVG(rating), -1) FROM reviews WHERE reviews.olid = sqlc.arg(olid) AND rating !=-1), (SELECT COUNT(reviews.olid) FROM reviews WHERE olid=sqlc.arg(olid) AND rating != -1))
ON CONFLICT(olid) DO UPDATE SET rating = excluded.rating;
