package db

import (
	"vidz/video"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(dbpath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&video.Video{})
	return db
}
