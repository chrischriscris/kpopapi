package model

import "gorm.io/gorm"

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
    Birthdate string
    Height float64
    Weight float64
}
