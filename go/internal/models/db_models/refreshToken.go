package db_models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	HashedToken 	 string `gorm:"unique;not null"`
	UserID   uint   `gorm:"not null"`
	User	 User   `gorm:"foreignKey:UserID;references:ID"`
	ExpiresAt time.Time `gorm:"not null"`
}

func (RefreshToken) TableName() string {
	return "RefreshToken"
}
