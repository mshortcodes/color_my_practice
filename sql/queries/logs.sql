-- name: CreateLog :one
INSERT INTO logs (id, date, color_depth, user_id)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetLogs :many
SELECT * FROM logs
ORDER BY date DESC;

-- name: GetLog :one
SELECT * FROM logs
WHERE id = $1;

-- name: DeleteLog :exec
DELETE FROM logs
WHERE id = $1;

-- name: ConfirmLogs :many
UPDATE logs
SET confirmed = true
WHERE user_id = $1
AND id = ANY($2::UUID[])
RETURNING *;