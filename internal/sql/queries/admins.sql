-- name: CreateAdmin :one
INSERT INTO admins (id, user_id)
VALUES (?, ?)
RETURNING *;

-- name: DeleteAdmin :exec
DELETE FROM admins
WHERE user_id = ?;

-- name: IsAdmin :one
SELECT * FROM admins
WHERE user_id = ?;