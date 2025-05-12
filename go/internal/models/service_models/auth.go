package service_models

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
	Token string
}

type UpdatePasswordRequest struct {
	UserID      uint
	OldPassword string
	NewPassword string
}
