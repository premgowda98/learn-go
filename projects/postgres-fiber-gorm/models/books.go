package models

import "gorm.io/gorm"

type Books struct {
	ID        uint   `gorm:"primary key;autoincrement" json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate()

	return err
}