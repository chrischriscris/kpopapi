package model

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	URL      string
	Metadata ImageMetadata
}

type ImageMetadata struct {
	gorm.Model
	ImageID uint
	Width   uint
	Height  uint
	Size    uint
	Format  string
    Landscape bool
}
