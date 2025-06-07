-- name: UpsertPosts :many
INSERT INTO
    posts
SELECT p.*
FROM UNNEST($1::public.posts[]) AS p
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
