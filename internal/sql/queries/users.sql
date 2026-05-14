-- name: CreateUser :one
INSERT INTO users (id, username, password)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserById :one
SELECT id, username, created_at, updated_at
FROM users
WHERE id = ?;

-- name: GetUserByName :one
SELECT *
FROM users
WHERE username = ?;

-- name: GetUserByRefreshToken :one
SELECT users.* FROM users
INNER JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = ?
    AND refresh_tokens.expires_at > CURRENT_TIMESTAMP
    AND refresh_tokens.revoked_at IS NULL;