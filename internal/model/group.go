package model

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
    gorm.Model
    Name string
    CompanyID uint
    Info *GroupInfo
    Type string // Either Boy Group, Girl Group or Co-Ed Group
    Members []*Idol `gorm:"many2many:group_idols;"`
}

type GroupInfo struct {
    gorm.Model
    GroupID uint
    DebutDate time.Time
}

