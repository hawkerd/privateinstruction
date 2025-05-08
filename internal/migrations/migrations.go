package migrations

import (
	"log"

	"github.com/hawkerd/privateinstruction/internal/models/db_models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// run migrations for all models
	err := db.AutoMigrate(
		&db_models.User{},
		&db_models.Class{},
		&db_models.ClassMember{},
	)
	if err != nil {
		return err
	}

	log.Println("Database migrated successfully")
	return nil
}
