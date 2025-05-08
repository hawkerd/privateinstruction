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
