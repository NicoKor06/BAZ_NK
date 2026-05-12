-- name: CreateBlog :one
INSERT INTO blogs (headline, body, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBlogByID :one
SELECT * FROM blogs WHERE blog_id = $1 LIMIT 1;

-- name: ListBlogs :many
SELECT * FROM blogs
ORDER BY created_at DESC
OFFSET $1 LIMIT $2;

-- name: CountBlogs :one
SELECT COUNT(*) FROM blogs;

-- name: UpdateBlog :one
UPDATE blogs SET
    headline = COALESCE($2, headline), -- Falls ein neuer Wert vorhanden ist, ersetze ihn. Ansonsten behalte den Alten.
    body = COALESCE($3, body),
    updated_at = NOW()
WHERE blog_id = $1 AND user_id = $4
RETURNING *;

-- name: DeleteBlog :exec
DELETE FROM blogs WHERE blog_id = $1 AND user_id = $2;

-- name: DeleteBlogsByUser :exec
DELETE FROM blogs WHERE user_id = $1;

-- name: DeleteBlogByID :exec
DELETE FROM blogs WHERE blog_id = $1;