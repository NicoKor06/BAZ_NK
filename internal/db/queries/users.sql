-- name: CreateUser :one
INSERT INTO users (
    username, firstname, lastname, email, password, birthday, role, last_online
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE user_id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users SET
    firstname = COALESCE($2, firstname),
    lastname = COALESCE($3, lastname),
    email = COALESCE($4, email),
    birthday = COALESCE($5, birthday),
    updated_at = NOW()
WHERE user_id = $1
RETURNING *;

-- name: UpdateLastOnline :exec
UPDATE users SET last_online = NOW() WHERE user_id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE user_id = $1;

-- name: GetUserPublic :one
SELECT user_id, username, firstname, lastname, role, last_online 
FROM users WHERE user_id = $1 LIMIT 1;