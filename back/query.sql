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
