-- name: CreateTLDR :one
INSERT INTO tldrs (id, title, content, user_id)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetTLDRsFromUser :many
SELECT * FROM tldrs
WHERE user_id = ?;

-- name: GetTLDRFromUser :one
SELECT * FROM tldrs
WHERE user_id = ?
    AND id = ?;

-- name: GetTLDRById :one
SELECT * FROM tldrs
WHERE id = ?;

-- name: GetAllTLDRs :many
SELECT * FROM tldrs;

-- name: UpdateTLDRTitle :one
UPDATE tldrs
SET title = ?
WHERE user_id = ?
    AND id = ?
RETURNING *;

-- name: DeleteTLDR :exec
DELETE FROM tldrs
WHERE user_id = ?
    AND id = ?;