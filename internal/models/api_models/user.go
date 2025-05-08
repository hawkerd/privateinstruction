package api_models

// read user
type ReadUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
