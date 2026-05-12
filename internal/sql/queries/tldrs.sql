-- name: CreateTLDR :one
INSERT INTO tldrs (id, title, content, user_id)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetTLDRsByUser :many
SELECT * FROM tldrs
WHERE user_id = ?;

-- name: GetTLDRById :one
SELECT * FROM tldrs
WHERE user_id = ?
    AND id = ?;

-- name: UpdateTLDRTitleById :exec
UPDATE tldrs
SET title = ?
WHERE user_id = ?
    AND id = ?;

-- name: DeleteTLDRById :exec
DELETE FROM tldrs
WHERE user_id = ?
    AND id = ?;