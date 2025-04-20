package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hawkerd/privateinstruction/pkg/auth"
	"github.com/hawkerd/privateinstruction/pkg/models"
	"gorm.io/gorm"
)

var db *gorm.DB

// set the database connection for the handlers
func SetDB(database *gorm.DB) {
	db = database
}

// sign up a new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// decode the request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the user input
	if user.Username == "" || user.Password == "" || user.Email == "" {
		http.Error(w, "Username, password, and email are required", http.StatusBadRequest)
		return
	}

	// check if the user already exists (email or username)
	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "User already exists with that email", http.StatusConflict)
		return
	}
	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		http.Error(w, "User already exists with that username", http.StatusConflict)
		return
	}

	// hash the passsword
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create the new user
	user.Password = hashedPassword
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
