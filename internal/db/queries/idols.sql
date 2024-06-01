-- name: GetIdol :one
SELECT * FROM idols
WHERE id = ? LIMIT 1;

-- name: ListIdols :many
SELECT * FROM idols
ORDER BY name;

-- name: CreateIdol :one
INSERT INTO idols (
    stage_name,
    name,
    gender,
    idol_info_id
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateIdol :exec
UPDATE idols
set name = ?
WHERE id = ?;

-- name: DeleteAuthor :exec
DELETE FROM idols
WHERE id = ?;
