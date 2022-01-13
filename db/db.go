package db

import (
	"fmt"
	"os"
	"vidz/video"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(dbpath string) *gorm.DB {
	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		fmt.Println("Creating database file: ", dbpath)
		_, err = os.Create(dbpath)
		if err != nil {
			panic("Could not create database file: " + err.Error())
		}
	}

	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&video.Video{})
	return db
}
