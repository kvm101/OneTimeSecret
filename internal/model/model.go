package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;not null"`
	Username *string    `gorm:"type:varchar(100);not null;unique"`
	Password *string    `gorm:"type:char(64);not null"`
	Messages *[]Message `gorm:"foreignKey:UserID"`
}

type Message struct {
	ID              *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;not null"`
	Text            *string    `gorm:"type:text;default:'';not null"`
	Timestamp       *time.Time `gorm:"type:timestamp;default:current_timestamp;not null"`
	ExpirationDate  *time.Time `gorm:"type:timestamp;default:null"`
	MessagePassword *string    `gorm:"type:char(64);default:''"`
	Times           *int       `gorm:"type:integer;default:null"`
	UserID          *uuid.UUID `gorm:"type:uuid;index"`
	User            *User      `gorm:"foreignKey:UserID"`
}

type AccountData struct {
	Username string
	Messages *[]Message
	IsAuth   bool
}

type MessageInfo struct {
	ID        *uuid.UUID
	Text      *string
	Times     *int
	Timestamp *time.Time
	Username  *string
}
