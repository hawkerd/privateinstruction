package api_models

// read user
type ReadUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// update user
type UpdateUserRequest struct {
	UserID   uint
	Username string
	Email    string
}
