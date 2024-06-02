-- name: GetIdolByName :one
SELECT * FROM idols
WHERE name = $1 or stage_name = $1
LIMIT 1;

-- name: ListIdols :many
SELECT * FROM idols;

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
