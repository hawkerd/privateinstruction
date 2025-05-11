package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/hawkerd/privateinstruction/docs"
	"github.com/hawkerd/privateinstruction/internal/db"
	"github.com/hawkerd/privateinstruction/internal/handlers"
	"github.com/hawkerd/privateinstruction/internal/middleware"
	"github.com/hawkerd/privateinstruction/internal/migrations"
	"github.com/hawkerd/privateinstruction/internal/services"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// establish db connection and run the migrations
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := migrations.Migrate(dbConn); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	authService := services.NewAuthService(dbConn)
	userService := services.NewUserService(dbConn)
	classService := services.NewClassService(dbConn)

	// create a router
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Post("/signup", handlers.SignUp(authService))
	r.Post("/signin", handlers.SignIn(authService))

	r.Group(func(r chi.Router) {
		r.Use(middleware.TokenAuthMiddleware)
		r.Get("/me", handlers.ReadUser(userService))
		r.Delete("/me", handlers.DeleteUser(userService))
		r.Put("/me", handlers.UpdateUser(userService))
		r.Put("/me/password", handlers.UpdatePassword(authService))

		r.Post("/class", handlers.CreateClass(classService))
		r.Delete("/class/{id}", handlers.DeleteClass(classService))
		r.Get("/class/{id}", handlers.ReadClass(classService))
		r.Put("/class/{id}", handlers.UpdateClass(classService))
		//r.Post("/class", handlers.CreateClass)
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
