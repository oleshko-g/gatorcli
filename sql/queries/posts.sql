-- name: InsertPost :one
INSERT INTO
    posts (
        id,
        created_at,
        updated_at,
        title,
        url,
        description,
        published_at,
        feed_id
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    )
ON CONFLICT (url) DO NOTHING
RETURNING
    *;

-- name: GetPostsForUser :many
WITH
    feed_ids_user_follows AS (
        SELECT feed_follows.feed_id
        FROM feed_follows
        WHERE
            feed_follows.user_id = $1
    )
SELECT posts.*
FROM posts
    JOIN feed_ids_user_follows USING (feed_id)
LIMIT $2;