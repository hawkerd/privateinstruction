package services

import (
	"errors"
	"strings"

	"github.com/hawkerd/privateinstruction/internal/auth"
	"github.com/hawkerd/privateinstruction/internal/models/db_models"
	"github.com/hawkerd/privateinstruction/internal/models/service_models"
	"gorm.io/gorm"
)

// define custom error messages
var (
	ErrUserExists          = errors.New("account already exists with that email or username")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrTokenGeneration     = errors.New("failed to generate token")
	ErrInternalServerError = errors.New("internal server error")
)

type AuthService struct {
	DB *gorm.DB
}

// create and return a new AuthService instance
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		DB: db,
	}
}

// register a user
func (s *AuthService) SignUp(req service_models.SignUpRequest) error {
	// normalize input
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// check if the user already exists (email or username)
	var existingUser db_models.User
	err := s.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error

	// if the user exists (no error), return
	if err == nil {
		return ErrUserExists
	}

	// if the error is not a record not found error, return it
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrInternalServerError
	}

	// hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return ErrInternalServerError
	}

	// create the new user
	var user db_models.User
	user.Username = req.Username
	user.Email = req.Email
	user.HashedPassword = hashedPassword

	// create the user in the database
	if err := s.DB.Create(&user).Error; err != nil {
		return ErrInternalServerError
	}

	return nil
}

// attempt to authenticate a user, and return a JWT token if successful
func (s *AuthService) SignIn(req service_models.SignInRequest) (service_models.SignInResponse, error) {
	// find the user by username or email
	var user db_models.User
	if req.Username != "" {
		if err := s.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
			return service_models.SignInResponse{}, ErrInvalidCredentials
		}
	} else {
		if err := s.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
			return service_models.SignInResponse{}, ErrInvalidCredentials
		}
	}

	// check the input password against the hashed password
	if !auth.CheckPassword(user.HashedPassword, req.Password) {
		return service_models.SignInResponse{}, ErrInvalidCredentials
	}

	// generate a JWT token
	token, err := auth.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return service_models.SignInResponse{}, ErrTokenGeneration
	}

	// create the response
	res := service_models.SignInResponse{
		Token: token,
	}

	return res, nil
}

// update the user's password
func (s *AuthService) UpdatePassword(req service_models.UpdatePasswordRequest) error {
	// find the user by ID
	var user db_models.User
	if err := s.DB.First(&user, req.UserID).Error; err != nil {
		return ErrInvalidCredentials
	}

	// check the input password against the hashed password
	if !auth.CheckPassword(user.HashedPassword, req.OldPassword) {
		return ErrInvalidCredentials
	}

	// hash the new password
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return ErrInternalServerError
	}

	// update the user's password in the database
	user.HashedPassword = hashedPassword
	if err := s.DB.Save(&user).Error; err != nil {
		return ErrInternalServerError
	}

	return nil
}
