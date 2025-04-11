package model

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() error {
	dsn := "host=localhost user=postgres password=admin dbname=postgres port=2345 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return err
	}

	DB = db

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return err
	}

	err = db.AutoMigrate(&User{}, &Message{})
	if err != nil {
		log.Println("AutoMigrate failed:", err)
		return err
	}

	log.Println("Database connected and migrated successfully.")
	return nil
}
