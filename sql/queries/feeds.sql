-- name: CreateFeed :one
INSERT INTO
    feeds (
        id,
        created_at,
        updated_at,
        name,
        url,
        user_id
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetFeedsUsers :many
SELECT
    sqlc.embed(feeds),
    sqlc.embed(users)
FROM feeds
    JOIN users ON feeds.user_id = users.id;

-- name: CreateFeedFollow :one 
WITH
    cte_inserted_feed_follow AS (
        INSERT INTO
            feed_follows (
                id,
                created_at,
                updated_at,
                user_id,
                feed_url
            )
        VALUES ($1, $2, $3, $4, $5)
    )
SELECT sqlc.embed(feed_follows), sqlc.embed(users), sqlc.embed(feeds)
FROM
    feed_follows
    JOIN users ON feed_follows.user_id = users.id
    JOIN feeds ON feed_follows.feed_url = feeds.url
    WHERE feeds.id = $1;


-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;