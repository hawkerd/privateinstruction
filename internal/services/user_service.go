package services

import (
	"errors"

	"github.com/hawkerd/privateinstruction/internal/models/db_models"
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
func (s *UserService) ReadUser(userID uint) (*db_models.User, error) {
	var user db_models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
