-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING *;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE id = $1;

-- name: GetFeedFollows :many
SELECT *
FROM feed_follows
WHERE user_id = $1;
