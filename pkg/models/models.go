package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}

type Class struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
}

type ClassMember struct {
	gorm.Model
	ClassID uint   `gorm:"not null"`
	Class   Class  `gorm:"foreignKey:ClassID"`
	UserID  uint   `gorm:"not null"`
	User    User   `gorm:"foreignKey:UserID"`
	Role    string `gorm:"not null"` // e.g., "student", "teacher"
}
