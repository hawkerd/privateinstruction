package db_models

import (
	"gorm.io/gorm"
)

type Class struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	CreatorID   uint
	CreatedBy   User `gorm:"foreignKey:CreatorID"`
}

func (Class) TableName() string {
	return "Class"
}
