-- name: CreateUser :one
INSERT INTO users (id, username, password)
VALUES (?, ?, ?)
RETURNING id, username, created_at, updated_at;

-- name: GetUserById :one
SELECT id, username, created_at, updated_at
FROM users
WHERE id = ?;

-- name: GetUserByName :one
SELECT id, username, created_at, updated_at
FROM users
WHERE username = ?;