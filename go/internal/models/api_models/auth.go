package api_models

// sign up
type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// sign in
type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
type SignInResponse struct {
	AccessToken string `json:"access_token"`
}

// update password
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// refresh token
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
}