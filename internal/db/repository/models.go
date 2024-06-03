// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package repository

import (
	"database/sql"
	"time"
)

type Company struct {
	ID           int32
	Name         string
	Country      string
	CreationDate sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Group struct {
	ID        int32
	Name      string
	Type      string
	DebutDate sql.NullTime
	CompanyID sql.NullInt32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GroupImage struct {
	ID        int32
	GroupID   int32
	ImageID   int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GroupMember struct {
	ID        int32
	GroupID   int32
	IdolID    int32
	SinceDate sql.NullTime
	UntilDate sql.NullTime
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Idol struct {
	ID        int32
	StageName string
	Name      sql.NullString
	Gender    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IdolImage struct {
	ID        int32
	IdolID    int32
	ImageID   int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IdolInfo struct {
	ID        int32
	IdolID    int32
	Birthdate sql.NullTime
	HeightCm  sql.NullFloat64
	WeightKg  sql.NullFloat64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Image struct {
	ID        int32
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ImageMetadatum struct {
	ID        int32
	ImageID   int32
	Width     int32
	Height    int32
	Landscape bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
