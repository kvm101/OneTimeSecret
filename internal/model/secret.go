package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username string    `gorm:"type:varchar(100);not null"`
	Password string    `gorm:"type:char(64);not null"`
	Messages []Message `gorm:"foreignKey:UserID"`
}

type Message struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Text            string    `gorm:"type:text"`
	Timestamp       time.Time `gorm:"type:timestamp;default:current_timestamp"`
	ExpirationDate  time.Time `gorm:"type:timestamp;default:null"`
	MessagePassword string    `gorm:"type:char(64);default:''"`
	UserID          uuid.UUID `gorm:"type:uuid;index"`
	User            User      `gorm:"foreignKey:UserID"`
}
