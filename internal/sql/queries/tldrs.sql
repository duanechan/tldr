-- name: CreateTLDR :one
INSERT INTO tldrs (id, title, content, user_id)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetTLDRsByUser :many
SELECT id, created_at, updated_at, title FROM tldrs
WHERE user_id = ?
    AND created_at < ?
ORDER BY created_at DESC
LIMIT ?;

-- name: GetTLDRByIDAndUser :one
SELECT * FROM tldrs
WHERE user_id = ?
    AND id = ?;

-- name: GetTLDRById :one
SELECT * FROM tldrs
WHERE id = ?;

-- name: GetTLDRs :many
SELECT id, created_at, updated_at, title FROM tldrs
WHERE created_at < ?
ORDER BY created_at DESC
LIMIT ?;

-- name: UpdateTLDRTitle :one
UPDATE tldrs
SET title = ?
WHERE user_id = ?
    AND id = ?
RETURNING *;

-- name: UpdateTLDRTitleById :one
UPDATE tldrs
SET title = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTLDR :exec
DELETE FROM tldrs
WHERE user_id = ?
    AND id = ?;

-- name: DeleteTLDRById :exec
DELETE FROM tldrs
WHERE id = ?;