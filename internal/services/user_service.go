package services

import (
	"errors"

	"github.com/hawkerd/privateinstruction/internal/models/db_models"
	"github.com/hawkerd/privateinstruction/internal/models/service_models"
	"gorm.io/gorm"
)

// define custom error messages
var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	DB *gorm.DB
}

// create and return a new AuthService instance
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

// read a user by ID
func (s *UserService) ReadUser(req service_models.ReadUserRequest) (*service_models.ReadUserResponse, error) {
	var user db_models.User
	if err := s.DB.First(&user, req.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// build the response
	response := service_models.ReadUserResponse{
		Username: user.Username,
		Email:    user.Email,
	}

	return &response, nil
}

// delete a user by ID
func (s *UserService) DeleteUser(req service_models.DeleteUserRequest) error {
	var user db_models.User
	if err := s.DB.First(&user, req.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if err := s.DB.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

// update a user by ID
func (s *UserService) UpdateUser(req service_models.UpdateUserRequest) error {
	var user db_models.User
	if err := s.DB.First(&user, req.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	user.Username = req.Username
	user.Email = req.Email

	if err := s.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
