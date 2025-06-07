-- name: UpsertPosts :many
INSERT INTO
    posts
SELECT p.*
FROM UNNEST($1::public.posts[]) AS p
ON CONFLICT (url) DO NOTHING
RETURNING
    *;