package db_models

import (
	"time"

	"gorm.io/gorm"
)

type JoinCode struct {
	gorm.Model
	Code         string `gorm:"not null"`
	ClassID      uint
	Class        Class     `gorm:"foreignKey:ClassID"`
	ExpirationDT time.Time `gorm:"not null"`
}

func (JoinCode) TableName() string {
	return "JoinCode"
}
