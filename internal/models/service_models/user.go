package service_models

type ReadUserRequest struct {
	UserID uint
}
type ReadUserResponse struct {
	Username string
	Email    string
}

type DeleteUserRequest struct {
	UserID uint
}

type UpdateUserRequest struct {
	UserID   uint
	Username string
	Email    string
}
