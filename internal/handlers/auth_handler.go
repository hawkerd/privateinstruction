package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hawkerd/privateinstruction/internal/models/api_models"
	"github.com/hawkerd/privateinstruction/internal/models/service_models"
	"github.com/hawkerd/privateinstruction/internal/services"
)

// sign up
func SignUp(authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decode the request body
		var req api_models.SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// input validation
		if req.Username == "" || req.Password == "" || req.Email == "" {
			http.Error(w, "username, password, and email are required", http.StatusBadRequest)
			return
		}

		// build the service request
		sreq := service_models.SignUpRequest{
			Username: req.Username,
			Password: req.Password,
			Email:    req.Email,
		}

		// call the service
		if err := authService.SignUp(sreq); err != nil {
			// return appropriate error message
			if errors.Is(err, services.ErrUserExists) {
				http.Error(w, err.Error(), http.StatusConflict)
			} else {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		// build the response
		res := api_models.SignUpResponse{
			Message: "user created successfully",
		}

		// encode the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// sign in
func SignIn(authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decode the request body
		var req api_models.SignInRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// input validation
		if (req.Username == "" && req.Email == "") || req.Password == "" {
			http.Error(w, "username/email and password are required", http.StatusBadRequest)
			return
		}

		// build the service request
		sreq := service_models.SignInRequest{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		// call the service
		sres, err := authService.SignIn(sreq)
		if err != nil {
			// return appropriate error message
			if errors.Is(err, services.ErrInvalidCredentials) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		// build the response
		res := api_models.SignInResponse{
			Token: sres.Token,
		}

		// encode the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// update password
func UpdatePassword(authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decode the request body
		var req api_models.UpdatePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// input validation
		if req.UserID == 0 || req.OldPassword == "" || req.NewPassword == "" {
			http.Error(w, "user_id, old_password, and new_password are required", http.StatusBadRequest)
			return
		}

		// build the service request
		sreq := service_models.UpdatePasswordRequest{
			UserID:      req.UserID,
			OldPassword: req.OldPassword,
			NewPassword: req.NewPassword,
		}

		// call the service
		if err := authService.UpdatePassword(sreq); err != nil {
			// return appropriate error message
			if errors.Is(err, services.ErrInvalidCredentials) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else if errors.Is(err, services.ErrUserNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
