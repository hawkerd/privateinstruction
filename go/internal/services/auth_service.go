package services

import (
	"errors"
	"strings"
	"time"

	"github.com/hawkerd/privateinstruction/internal/auth"
	"github.com/hawkerd/privateinstruction/internal/models/db_models"
	"github.com/hawkerd/privateinstruction/internal/models/service_models"
	"gorm.io/gorm"
)

// define custom error messages
var (
	ErrUserExists          = errors.New("Account already exists with this email or username")
	ErrInvalidCredentials  = errors.New("Invalid credentials")
	ErrTokenGeneration     = errors.New("failed to generate token")
	ErrInternalServerError = errors.New("Something went wrong")
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
	user := db_models.User{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

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
	accessToken, err := auth.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return service_models.SignInResponse{}, ErrTokenGeneration
	}

	// generate a refresh token and hash it
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return service_models.SignInResponse{}, ErrTokenGeneration
	}
	hashedToken, err := auth.HashPassword(refreshToken)
	if err != nil {
		return service_models.SignInResponse{}, ErrInternalServerError
	}
	expiration := auth.RefreshTokenExpiration()

	// store the refresh token in the database
	refreshTokenRecord := db_models.RefreshToken{
		UserID:      user.ID,
		HashedToken: hashedToken,
		ExpiresAt:   expiration,
	}
	if err := s.DB.Create(&refreshTokenRecord).Error; err != nil {
		return service_models.SignInResponse{}, ErrInternalServerError
	}

	// create the response
	res := service_models.SignInResponse{
		AccessToken:            accessToken,
		RefreshToken:           refreshToken,
		RefreshTokenExpiration: expiration,
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

// generate a new access token
func (s *AuthService) RefreshAccessToken(req service_models.RefreshTokenRequest) (service_models.RefreshTokenResponse, error) {
	var refreshTokens []db_models.RefreshToken
	if err := s.DB.Where("user_id = ?", req.UserID).Find(&refreshTokens).Error; err != nil {
		return service_models.RefreshTokenResponse{}, ErrInternalServerError
	}

	// match the passed refresh token with the stored hashed tokens
	var refreshToken *db_models.RefreshToken
	for i, token := range refreshTokens {
		if auth.CheckPassword(token.HashedToken, req.RefreshToken) {
			refreshToken = &refreshTokens[i]
			break
		}
	}
	if refreshToken == nil {
		return service_models.RefreshTokenResponse{}, ErrInvalidCredentials
	}

	// check if the refresh token is expired
	if time.Now().After(refreshToken.ExpiresAt) {
		return service_models.RefreshTokenResponse{}, ErrInvalidCredentials
	}

	// query the user by ID
	var user db_models.User
	if err := s.DB.First(&user, refreshToken.UserID).Error; err != nil {
		return service_models.RefreshTokenResponse{}, ErrInvalidCredentials
	}
	// check if the user exists
	if user.ID == 0 {
		return service_models.RefreshTokenResponse{}, ErrInvalidCredentials
	}

	// generate a new access token
	accessToken, err := auth.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return service_models.RefreshTokenResponse{}, ErrTokenGeneration
	}

	// generate a new refresh token
	newRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return service_models.RefreshTokenResponse{}, ErrTokenGeneration
	}
	// hash the new refresh token
	newHashedToken, err := auth.HashPassword(newRefreshToken)
	if err != nil {
		return service_models.RefreshTokenResponse{}, ErrInternalServerError
	}
	// update the new refresh token in the database
	refreshToken.HashedToken = newHashedToken
	refreshToken.ExpiresAt = auth.RefreshTokenExpiration()
	if err := s.DB.Save(&refreshToken).Error; err != nil {
		return service_models.RefreshTokenResponse{}, ErrInternalServerError
	}

	// create the response
	res := service_models.RefreshTokenResponse{
		AccessToken:            accessToken,
		RefreshToken:           newRefreshToken,
		RefreshTokenExpiration: refreshToken.ExpiresAt,
	}

	return res, nil
}
