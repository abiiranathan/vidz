package db

import (
	"vidz/video"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("videos.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&video.Video{})

	return db
}
