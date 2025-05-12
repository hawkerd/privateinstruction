package service_models

import "time"

type CreateClassRequest struct {
	Name        string
	Description string
	UserID      uint
}

type DeleteClassRequest struct {
	ClassID uint
	UserID  uint
}

type ReadClassRequest struct {
	ClassID uint
	UserID  uint
}
type ReadClassResponse struct {
	Name        string
	Description string
	CreatedAt   string
	CreatedBy   string
}

type UpdateClassRequest struct {
	Name        string
	Description string
	UserID      uint
	ClassID     uint
}

type GenerateJoinCodeRequest struct {
	ClassID uint
	UserID  uint
}

type GenerateJoinCodeResponse struct {
	Code         string
	ExpirationDT time.Time
}

type JoinClassRequest struct {
	JoinCode string
	UserID   uint
}
