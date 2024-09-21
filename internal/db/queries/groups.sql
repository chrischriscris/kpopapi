-- name: GetGroupByName :one
SELECT * FROM groups
WHERE name = $1
LIMIT 1;

-- name: ListGroups :many
SELECT * FROM groups;

-- name: CreateGroupMinimal :one
INSERT INTO groups (name, type)
VALUES ($1, $2)
RETURNING *;

-- name: CreateGroupMinimalWithDebut :one
INSERT INTO groups (name, type, debut_date)
VALUES ($1, $2, $3)
RETURNING *;
