-- name: CreateFeed :one
INSERT INTO feeds (id, created_At, updated_at, name, url, user_id, last_fetched_at)
VALUES ($1, NOW(), NOW(), $2, $3, $4, NULL)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at IS NULL DESC
LIMIT $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at=NOW(),
updated_at=NOW()
WHERE url=$1;
