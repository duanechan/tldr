-- name: CreateUser :one
INSERT INTO users (id, username, password)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetUsers :many
SELECT id, created_at, updated_at, username
FROM users
WHERE created_at < ?
    OR (created_at = ? AND id < ?)
ORDER BY created_at DESC
LIMIT ?;

-- name: GetUserById :one
SELECT id, created_at, updated_at, username
FROM users
WHERE id = ?;

-- name: GetUserCredentialsByUsername :one
SELECT id, username, password
FROM users
WHERE username = ?;

-- name: GetUserByRefreshToken :one
SELECT users.*
FROM users
INNER JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = ?
    AND refresh_tokens.expires_at > CURRENT_TIMESTAMP
    AND refresh_tokens.revoked_at IS NULL;

-- name: UpdateUsername :one
UPDATE users
SET username = ?
WHERE id = ?
RETURNING *;

-- name: UpdatePassword :one
UPDATE users
SET password = ?
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: DeleteUsers :exec
DELETE FROM users
WHERE id IN (sqlc.slice('ids'));