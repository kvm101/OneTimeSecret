package model

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() error {
	dns := "host=localhost user=postgres password=admin dbname=postgres port=2345 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}

	DB = db

	return db.AutoMigrate(&User{}, &Message{})
}
