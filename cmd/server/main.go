package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hawkerd/privateinstruction/pkg/config"
	"github.com/hawkerd/privateinstruction/pkg/handlers"
	"github.com/hawkerd/privateinstruction/pkg/models"
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

	err = db.AutoMigrate(
		&models.User{},
		&models.Class{},
		&models.ClassMember{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database schema: %v", err)
	}

	handlers.SetDB(db)

	// create a router
	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Request:", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	r.Post("/signup", handlers.SignUp)
	r.Post("/signin", handlers.SignIn)
	r.Get("/me", handlers.GetUserInfo)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
