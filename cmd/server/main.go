package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hawkerd/privateinstruction/pkg/config"
	"github.com/hawkerd/privateinstruction/pkg/handlers"
	"github.com/hawkerd/privateinstruction/pkg/middleware"
	"github.com/hawkerd/privateinstruction/pkg/models"
	"github.com/rs/cors"
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

	r.Post("/signup", handlers.SignUp)
	r.Post("/signin", handlers.SignIn)

	r.Group(func(r chi.Router) {
		r.Use(middleware.TokenAuthMiddleware)
		r.Get("/me", handlers.ReadUser)
		r.Post("/class", handlers.CreateClass)
		//r.Get("/classes", handlers.GetClasses)
	})

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Frontend URL
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})
	handler := c.Handler(r)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", handler)
}
