-- name: InsertISBN :exec
INSERT INTO isbns (isbn, olid) VALUES (?, ?);

-- name: GetReviewsByOLID :many
SELECT * FROM reviews WHERE olid = ?;

-- name: InsertReview :exec
INSERT INTO reviews ( olid, source, external_id, username, rating, text) values (?, ?, ?, ?, ?, ?);

-- name: InsertWork :exec
INSERT INTO works (olid, title, author, description, published_year) values (?, ?, ?, ?, ?)
