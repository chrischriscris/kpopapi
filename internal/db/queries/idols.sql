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
  gender
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: CreateIdolWithGroupMinimal :one
INSERT INTO idols (
  stage_name,
  gender
) VALUES (
  $1, $2
)
RETURNING *;

-- name: AddMemberToGroup :one
INSERT INTO group_members (
  group_id,
  idol_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetIdolsByNameLike :many
SELECT * FROM idols
WHERE stage_name ILIKE '%' || $1 || '%'
OR name ILIKE '%' || $1 || '%';
