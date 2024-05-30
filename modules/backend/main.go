package main

import (
	"backend/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
  db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

  // Migrate the schema
  db.AutoMigrate(&model.Idol{}, &model.Group{}, &model.Company{}, &model.IdolInfo{}, &model.GroupInfo{})
}
