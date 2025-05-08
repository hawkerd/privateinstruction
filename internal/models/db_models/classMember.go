package db_models

import (
	"gorm.io/gorm"
)

type ClassMember struct {
	gorm.Model
	ClassID uint   `gorm:"not null"`
	Class   Class  `gorm:"foreignKey:ClassID"`
	UserID  uint   `gorm:"not null"`
	User    User   `gorm:"foreignKey:UserID"`
	Role    string `gorm:"not null"` // e.g., "student", "teacher"
}
