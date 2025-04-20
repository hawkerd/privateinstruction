package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hawkerd/privateinstruction/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// load environment variables
	config.LoadEnv()

	// connect to the database
	dsn := config.GetDatabaseURL()
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// create a router
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
