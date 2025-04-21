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

// sign up
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

// sign in
func SignIn(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	// decode the request body into the credentials struct
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// make sure the proper fields are provided
	if (credentials.Username == "" && credentials.Email == "") || credentials.Password == "" {
		http.Error(w, "Username/email and password are required", http.StatusBadRequest)
		return
	}

	// find the user by username or email
	var user models.User
	if credentials.Username != "" {
		if err := db.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
			http.Error(w, "User not found with that username", http.StatusUnauthorized)
			return
		}
	} else {
		if err := db.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
			http.Error(w, "User not found with that email address", http.StatusUnauthorized)
			return
		}
	}

	// check the password
	if !auth.CheckPassword(user.Password, credentials.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// generate a JWT token
	token, err := auth.GenerateJWT(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// respond with the token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// get user info
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	// make sure the user is authenticated
	userID := r.Context().Value("userID")
	userIDInt, ok := userID.(uint)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := db.First(&user, userIDInt).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// respond with the user info
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
