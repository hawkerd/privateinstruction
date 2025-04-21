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

func isUserAuthenticated(r *http.Request) (uint, bool) {
	// extract the token from the request header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return 0, false
	}

	// parse the token
	claims, err := auth.ParseJWT(tokenString)
	if err != nil {
		return 0, false
	}

	// extract the user ID from the claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, false
	}

	return uint(userID), true
}

// get user info
func ReadUser(w http.ResponseWriter, r *http.Request) {
	// make sure the user is authenticated
	userID, ok := isUserAuthenticated(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// respond with the user info
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// create a class
func CreateClass(w http.ResponseWriter, r *http.Request) {
	// make sure the user is authenticated
	userID, ok := isUserAuthenticated(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// decode the request body into the class struct
	var class models.Class
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the class input
	if class.Name == "" {
		class.Name = "Unnamed Class"
	}

	// create the new class
	if err := db.Create(&class).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create a class member entry for the user
	classMember := models.ClassMember{
		ClassID: class.ID,
		UserID:  userID,
		Role:    "instructor",
	}

	if err := db.Create(&classMember).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(class)
}

//func GetClasses(w http.ResponseWriter, r *http.Request) {
//	// make sure the user is authenticated
//	userID, ok := isUserAuthenticated(r)
//	if !ok {
//		http.Error(w, "Unauthorized", http.StatusUnauthorized)
//		return
//	}
//
//	var classes []models.Class
//	if err := db.Model(&models.ClassMember{}).
//		Joins("JOIN classes ON class_members.class_id = classes.id").
//		Where("class_members.user_id = ?", userID).
//		Select("classes.*").
//		Scan(&classes).Error; err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(classes)
//}
