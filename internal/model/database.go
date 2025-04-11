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
		return err // Повертаємо помилку, якщо не вдалося підключитися до БД
	}

	DB = db

	// Міграція для таблиці User
	if !db.Migrator().HasTable(&User{}) {
		err = db.AutoMigrate(&User{})
		if err != nil {
			log.Fatal("Failed to migrate User:", err)
			return err // Повертаємо помилку, якщо міграція не вдалася
		}
	}

	// Міграція для таблиці Message
	if !db.Migrator().HasTable(&Message{}) {
		err = db.AutoMigrate(&Message{})
		if err != nil {
			log.Fatal("Failed to migrate Message:", err)
			return err // Повертаємо помилку, якщо міграція не вдалася
		}
	}

	// Якщо все успішно
	return nil
}
