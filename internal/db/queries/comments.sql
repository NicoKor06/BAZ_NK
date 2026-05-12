-- name: CreateComment :one
INSERT INTO comments (body, blog_id, user_id) 
VALUES ($1, $2, $3) 
RETURNING *;

-- name: GetCommentByID :one
SELECT * FROM comments WHERE comment_id = $1 LIMIT 1;

-- name: ListCommentsByBlogID :many
SELECT * FROM comments 
WHERE blog_id = $1 
ORDER BY created_at DESC 
OFFSET $2 LIMIT $3;

-- name: CountCommentsByBlogID :one
SELECT COUNT(*) FROM comments WHERE blog_id = $1;

-- name: UpdateComment :one
UPDATE comments SET
    body = $2,
    updated_at = NOW()
WHERE comment_id = $1 AND user_id = $3
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments WHERE comment_id = $1 AND user_id = $2; -- Der User selbst

-- name: DeleteCommentsByBlog :exec
DELETE FROM comments WHERE blog_id = $1;

-- name: DeleteCommentsByUser :exec
DELETE FROM comments WHERE user_id = $1;

-- name: DeleteCommentByID :exec
DELETE FROM comments WHERE comment_id = $1;