// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: idols.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addMemberToGroup = `-- name: AddMemberToGroup :one
INSERT INTO group_members (
  group_id,
  idol_id
) VALUES (
  $1, $2
)
RETURNING id, group_id, idol_id, since_date, until_date, created_at, updated_at
`

type AddMemberToGroupParams struct {
	GroupID int32
	IdolID  int32
}

func (q *Queries) AddMemberToGroup(ctx context.Context, arg AddMemberToGroupParams) (GroupMember, error) {
	row := q.db.QueryRow(ctx, addMemberToGroup, arg.GroupID, arg.IdolID)
	var i GroupMember
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.IdolID,
		&i.SinceDate,
		&i.UntilDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createIdol = `-- name: CreateIdol :one
INSERT INTO idols (
  stage_name,
  name,
  gender
) VALUES (
  $1, $2, $3
)
RETURNING id, stage_name, name, gender, created_at, updated_at
`

type CreateIdolParams struct {
	StageName string
	Name      pgtype.Text
	Gender    string
}

func (q *Queries) CreateIdol(ctx context.Context, arg CreateIdolParams) (Idol, error) {
	row := q.db.QueryRow(ctx, createIdol, arg.StageName, arg.Name, arg.Gender)
	var i Idol
	err := row.Scan(
		&i.ID,
		&i.StageName,
		&i.Name,
		&i.Gender,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createIdolWithGroupMinimal = `-- name: CreateIdolWithGroupMinimal :one
INSERT INTO idols (
  stage_name,
  gender
) VALUES (
  $1, $2
)
RETURNING id, stage_name, name, gender, created_at, updated_at
`

type CreateIdolWithGroupMinimalParams struct {
	StageName string
	Gender    string
}

func (q *Queries) CreateIdolWithGroupMinimal(ctx context.Context, arg CreateIdolWithGroupMinimalParams) (Idol, error) {
	row := q.db.QueryRow(ctx, createIdolWithGroupMinimal, arg.StageName, arg.Gender)
	var i Idol
	err := row.Scan(
		&i.ID,
		&i.StageName,
		&i.Name,
		&i.Gender,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getIdolByName = `-- name: GetIdolByName :one
SELECT id, stage_name, name, gender, created_at, updated_at FROM idols
WHERE name = $1 or stage_name = $1
LIMIT 1
`

func (q *Queries) GetIdolByName(ctx context.Context, name pgtype.Text) (Idol, error) {
	row := q.db.QueryRow(ctx, getIdolByName, name)
	var i Idol
	err := row.Scan(
		&i.ID,
		&i.StageName,
		&i.Name,
		&i.Gender,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getIdolsByNameLike = `-- name: GetIdolsByNameLike :many
SELECT id, stage_name, name, gender, created_at, updated_at FROM idols
WHERE stage_name ILIKE '%' || $1 || '%'
OR name ILIKE '%' || $1 || '%'
`

func (q *Queries) GetIdolsByNameLike(ctx context.Context, dollar_1 pgtype.Text) ([]Idol, error) {
	rows, err := q.db.Query(ctx, getIdolsByNameLike, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Idol
	for rows.Next() {
		var i Idol
		if err := rows.Scan(
			&i.ID,
			&i.StageName,
			&i.Name,
			&i.Gender,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listIdols = `-- name: ListIdols :many
SELECT id, stage_name, name, gender, created_at, updated_at FROM idols
`

func (q *Queries) ListIdols(ctx context.Context) ([]Idol, error) {
	rows, err := q.db.Query(ctx, listIdols)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Idol
	for rows.Next() {
		var i Idol
		if err := rows.Scan(
			&i.ID,
			&i.StageName,
			&i.Name,
			&i.Gender,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
