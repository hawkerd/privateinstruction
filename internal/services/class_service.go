package services

import (
	"errors"

	"github.com/hawkerd/privateinstruction/internal/models/db_models"
	"github.com/hawkerd/privateinstruction/internal/models/service_models"
	"gorm.io/gorm"
)

// define custom error messages
var (
	ErrClassNotFound = errors.New("class not found")
	ErrUnauthorized  = errors.New("unauthorized")
)

type ClassService struct {
	DB *gorm.DB
}

// create and return a new ClassService instance
func NewClassService(db *gorm.DB) *ClassService {
	return &ClassService{
		DB: db,
	}
}

// create a new class
func (s *ClassService) CreateClass(req service_models.CreateClassRequest) error {
	// create a new class
	class := db_models.Class{
		Name:        req.Name,
		Description: req.Description,
		CreatorID:   req.UserID,
	}
	if err := s.DB.Create(&class).Error; err != nil {
		return ErrInternalServerError
	}

	// add the creator of the class as a member (admin)
	classMember := db_models.ClassMember{
		ClassID: class.ID,
		UserID:  req.UserID,
		Role:    "admin",
	}
	if err := s.DB.Create(&classMember).Error; err != nil {
		return ErrInternalServerError
	}

	return nil
}

// delete a class
func (s *ClassService) DeleteClass(req service_models.DeleteClassRequest) error {
	// find the class
	var class db_models.Class
	if err := s.DB.First(&class, req.ClassID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrClassNotFound
		}
		return err
	}

	// find the class member
	var classMember db_models.ClassMember
	if err := s.DB.Where("class_id = ? AND user_id = ?", req.ClassID, req.UserID).First(&classMember).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUnauthorized
		}
	}

	// make sure the user is an admin
	if classMember.Role != "admin" {
		return ErrUnauthorized
	}

	// delete the class
	if err := s.DB.Delete(&class).Error; err != nil {
		return err
	}

	// delete all class members
	if err := s.DB.Where("class_id = ?", req.ClassID).Delete(&db_models.ClassMember{}).Error; err != nil {
		return err
	}

	return nil
}

// read a class
func (s *ClassService) ReadClass(req service_models.ReadClassRequest) (service_models.ReadClassResponse, error) {
	// find the class
	var class db_models.Class
	if err := s.DB.First(&class, req.ClassID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_models.ReadClassResponse{}, ErrClassNotFound
		}
		return service_models.ReadClassResponse{}, err
	}

	// find the class member
	var classMember db_models.ClassMember
	if err := s.DB.Where("class_id = ? AND user_id = ?", req.ClassID, req.UserID).First(&classMember).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_models.ReadClassResponse{}, ErrUnauthorized
		}
	}

	// find the user who create the class
	var creator db_models.User
	if err := s.DB.First(&creator, req.UserID).Error; err != nil {
		creator = db_models.User{}
	}

	// build the response
	resp := service_models.ReadClassResponse{
		Name:        class.Name,
		Description: class.Description,
		CreatedAt:   class.CreatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:   creator.Username,
	}

	return resp, nil
}

// update a class
func (s *ClassService) UpdateClass(req service_models.UpdateClassRequest) error {
	// find the class
	var class db_models.Class
	if err := s.DB.First(&class, req.ClassID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrClassNotFound
		}
		return err
	}

	// find the class member
	var classMember db_models.ClassMember
	if err := s.DB.Where("class_id = ? AND user_id = ?", req.ClassID, req.UserID).First(&classMember).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUnauthorized
		}
	}

	// make sure the user is an admin
	if classMember.Role != "admin" {
		return ErrUnauthorized
	}

	// update the class
	class.Name = req.Name
	class.Description = req.Description

	if err := s.DB.Save(&class).Error; err != nil {
		return err
	}

	return nil
}
