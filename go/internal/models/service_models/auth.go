package service_models

import (
	"time"
)

type SignUpRequest struct {
	Username string
	Password string
	Email    string
}

type SignInRequest struct {
	Username string
	Email    string
	Password string
}

type SignInResponse struct {
	AccessToken string
	RefreshToken string
	RefreshTokenExpiration time.Time
}

type UpdatePasswordRequest struct {
	UserID      uint
	OldPassword string
	NewPassword string
}

type RefreshTokenRequest struct {
	RefreshToken string
	UserID	  uint
}

type RefreshTokenResponse struct {
	AccessToken string
	RefreshToken string
	RefreshTokenExpiration time.Time
}