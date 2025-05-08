package db

import (
	"log"

	"github.com/hawkerd/privateinstruction/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// establish a connection to the PostgreSQL database
func ConnectDB() (*gorm.DB, error) {
	// load config
	config.LoadEnv()

	// get the database URL
	dsn := config.GetDatabaseURL()

	// connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database")
	return db, nil
}
