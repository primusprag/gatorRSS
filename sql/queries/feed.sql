-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ArchiveFeed :one
INSERT INTO feedsArchive (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ListFeeds :many
SELECT *
FROM feeds;

-- name: CreateFeedFollow :many
WITH inserted_feed_follows AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
    )
    RETURNING *
)
SELECT 
    inserted_feed_follows.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follows
INNER JOIN users ON users.id = user_id
INNER JOIN feeds ON feeds.id = feed_id;

-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: GetFeedByName :one
SELECT *
FROM feeds
WHERE name = $1;

-- name: GetArchivedFeedByName :one
SELECT *
FROM feedsArchive
WHERE name = $1;

-- name: DeleteFeed :one
DELETE FROM feeds
WHERE id = $1
RETURNING *;

-- name: GetFeedFollowsForUser :many
SELECT 
    feed_follows.*, 
    users.name AS user_name,
    feeds.name AS feed_name
FROM feed_follows
JOIN users ON users.id = user_id
JOIN feeds ON feeds.id = feed_id
WHERE feed_follows.user_id = $1;


-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;


-- name: DeleteArchivedFeed :one
DELETE FROM feedsArchive
WHERE id = $1
RETURNING *;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $2, updated_at = $2
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
Limit 1;