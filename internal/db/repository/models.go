// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Company struct {
	ID           int32
	Name         string
	Country      string
	CreationDate pgtype.Date
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
}

type Group struct {
	ID        int32
	Name      string
	Type      string
	DebutDate pgtype.Date
	CompanyID pgtype.Int4
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type GroupImage struct {
	ID        int32
	GroupID   int32
	ImageID   int32
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type GroupMember struct {
	ID        int32
	GroupID   int32
	IdolID    int32
	SinceDate pgtype.Date
	UntilDate pgtype.Date
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type Idol struct {
	ID        int32
	StageName string
	Name      pgtype.Text
	Gender    string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type IdolGroup struct {
	IdolName  string
	GroupName string
	GroupType string
}

type IdolImage struct {
	ID        int32
	IdolID    int32
	ImageID   int32
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type IdolInfo struct {
	ID        int32
	IdolID    int32
	Birthdate pgtype.Date
	HeightCm  pgtype.Float8
	WeightKg  pgtype.Float8
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type Image struct {
	ID        int32
	Url       string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type ImageMetadatum struct {
	ID        int32
	ImageID   int32
	Width     int32
	Height    int32
	Landscape bool
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}
