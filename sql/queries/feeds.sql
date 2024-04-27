-- name: CreateFeed :one
INSERT INTO feeds (id, created_At, updated_at, name, url, user_id, last_fetched_at)
VALUES ($1, NOW(), NOW(), $2, $3, $4, NULL)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;
