package config

import (
	"log"
	"one_time_secret/internal/model"

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

	return db.AutoMigrate(&model.User{}, &model.Message{})
}

func CleanInappropriateDB() {
	message := model.Message{}

	DB.Where("Text IS NULL").Delete(&message)
}
