-- name: InsertISBN :exec
INSERT INTO isbns (isbn, olid) VALUES (?, ?);

-- name: GetReviewsByOLID :many
SELECT * FROM reviews WHERE olid = ?;
