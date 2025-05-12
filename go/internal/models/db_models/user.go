package db_models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"unique;not null"`
	HashedPassword string `gorm:"not null"`
	Email          string `gorm:"unique;not null"`
}

func (User) TableName() string {
	return "User"
}
