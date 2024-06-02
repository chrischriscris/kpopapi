-- name: GetGroupByName :one
SELECT * FROM groups
WHERE name = $1
LIMIT 1;

-- name: ListGroups :many
SELECT * FROM groups;
