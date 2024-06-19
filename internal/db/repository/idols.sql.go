// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
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
SELECT i.id, i.stage_name, i.name, i.gender, i.created_at, i.updated_at, g.id, g.name, g.type, g.debut_date, g.company_id, g.created_at, g.updated_at
FROM idols i
LEFT JOIN group_members gm
ON i.id = gm.idol_id
LEFT JOIN groups g
ON gm.group_id = g.id
WHERE i.stage_name ILIKE '%' || $1 || '%'
OR i.name ILIKE '%' || $1 || '%'
`

type GetIdolsByNameLikeRow struct {
	ID          int32
	StageName   string
	Name        pgtype.Text
	Gender      string
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	ID_2        pgtype.Int4
	Name_2      pgtype.Text
	Type        pgtype.Text
	DebutDate   pgtype.Date
	CompanyID   pgtype.Int4
	CreatedAt_2 pgtype.Timestamp
	UpdatedAt_2 pgtype.Timestamp
}

func (q *Queries) GetIdolsByNameLike(ctx context.Context, dollar_1 pgtype.Text) ([]GetIdolsByNameLikeRow, error) {
	rows, err := q.db.Query(ctx, getIdolsByNameLike, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetIdolsByNameLikeRow
	for rows.Next() {
		var i GetIdolsByNameLikeRow
		if err := rows.Scan(
			&i.ID,
			&i.StageName,
			&i.Name,
			&i.Gender,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.Name_2,
			&i.Type,
			&i.DebutDate,
			&i.CompanyID,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
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
