package api_models

// sign up
type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
type SignUpResponse struct {
	Message string `json:"message"`
}

// sign in
type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
type SignInResponse struct {
	Token string `json:"token"`
}

// update password
type UpdatePasswordRequest struct {
	UserID      uint   `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
