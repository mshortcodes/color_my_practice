-- name: CreateLog :one
INSERT INTO logs (id, date, color_depth, user_id)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3
)
RETURNING *;