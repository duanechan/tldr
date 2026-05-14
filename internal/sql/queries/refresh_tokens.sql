-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (id, token, user_id, expires_at)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = CURRENT_TIMESTAMP, revoked_at = CURRENT_TIMESTAMP
WHERE token = ?;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = ?;