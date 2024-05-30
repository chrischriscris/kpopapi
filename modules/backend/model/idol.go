package model

import (
	"time"

	"gorm.io/gorm"
)

type Idol struct {
    gorm.Model
    Name string
    Gender string // Either Male, Female or Other
    Info *IdolInfo
    Groups []*Group `gorm:"many2many:group_idols;"`
}

type IdolInfo struct {
    gorm.Model
    IdolID uint
    Birthdate time.Time
    Height float64 // Centimeters
    Weight float64 // Kilograms
}
