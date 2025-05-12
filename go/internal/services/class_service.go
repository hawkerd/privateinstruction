package services

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"

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

// helper function to generate a random string
func RandomString(length int) string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err) // or handle error gracefully
		}
		b[i] = charset[num.Int64()]
	}
	return string(b)
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

// generate a join code for a class
func (s *ClassService) GenerateJoinCode(req service_models.GenerateJoinCodeRequest) (service_models.GenerateJoinCodeResponse, error) {
	// find the class
	var class db_models.Class
	if err := s.DB.First(&class, req.ClassID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_models.GenerateJoinCodeResponse{}, ErrClassNotFound
		}
		return service_models.GenerateJoinCodeResponse{}, err
	}

	// find the class member
	var classMember db_models.ClassMember
	if err := s.DB.Where("class_id = ? AND user_id = ?", req.ClassID, req.UserID).First(&classMember).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_models.GenerateJoinCodeResponse{}, ErrUnauthorized
		}
	}
	// make sure the user is an admin
	if classMember.Role != "admin" {
		return service_models.GenerateJoinCodeResponse{}, ErrUnauthorized
	}

	// remove any existing join codes for the class
	if err := s.DB.Where("class_id = ?", req.ClassID).Delete(&db_models.JoinCode{}).Error; err != nil {
		return service_models.GenerateJoinCodeResponse{}, err
	}

	// generate a join code
	joinCode := RandomString(8)
	expirationDT := time.Now().Add(24 * time.Hour)
	if err := s.DB.Create(&db_models.JoinCode{
		Code:         joinCode,
		ClassID:      req.ClassID,
		ExpirationDT: expirationDT,
	}).Error; err != nil {
		return service_models.GenerateJoinCodeResponse{}, err
	}

	// generate the response
	resp := service_models.GenerateJoinCodeResponse{
		Code:         joinCode,
		ExpirationDT: expirationDT,
	}

	return resp, nil
}

// join a class using a join code
func (s *ClassService) JoinClass(req service_models.JoinClassRequest) error {
	// find the join code
	var joinCode db_models.JoinCode
	if err := s.DB.Where("code = ?", req.JoinCode).First(&joinCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrClassNotFound
		}
		return err
	}

	// check if the join code is expired
	if time.Now().After(joinCode.ExpirationDT) {
		return ErrClassNotFound
	}

	// find the class
	var class db_models.Class
	if err := s.DB.First(&class, joinCode.ClassID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrClassNotFound
		}
		return err
	}

	// add the user to the class as a member (user)
	classMember := db_models.ClassMember{
		ClassID: class.ID,
		UserID:  req.UserID,
		Role:    "user",
	}
	if err := s.DB.Create(&classMember).Error; err != nil {
		return err
	}

	return nil
}
