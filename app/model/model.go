package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Image{})
	db.AutoMigrate(&Variant{})
	return db
}
