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
                feed_id
            )
        VALUES ($1, $2, $3, $4, $5)
        RETURNING *
    )
SELECT cte_inserted_feed_follow.*, sqlc.embed(users), sqlc.embed(feeds)
FROM
    cte_inserted_feed_follow
    JOIN users ON cte_inserted_feed_follow.user_id = users.id
    JOIN feeds ON cte_inserted_feed_follow.feed_id = feeds.id;


-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeedFollowUser :many
SELECT sqlc.embed(feed_follows), sqlc.embed(users), sqlc.embed(feeds) 
FROM users
    JOIN feed_follows ON users.id = feed_follows.user_id
    JOIN feeds on feed_follows.feed_id = feeds.id
WHERE
    users.id = $1;

-- name: DeleteFeedFollowUser :one
DELETE FROM feed_follows WHERE user_id = $1 and feed_id = $2 RETURNING *;