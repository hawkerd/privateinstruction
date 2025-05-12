package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hawkerd/privateinstruction/internal/models/api_models"
	"github.com/hawkerd/privateinstruction/internal/models/service_models"
	"github.com/hawkerd/privateinstruction/internal/services"
)

const userIDKey = "userID"

//	@Summary		ReadUser
//	@Description	Get user info
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"Bearer token"
//	@Router			/me [get]
//	@Security		Bearer
//	@Tags			User
func ReadUser(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the user ID from the request context
		userID, ok := r.Context().Value(userIDKey).(uint)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// build the service request
		sreq := service_models.ReadUserRequest{
			UserID: userID,
		}

		// call the service to get user info
		sres, err := userService.ReadUser(sreq)
		if err != nil {
			if errors.Is(err, services.ErrUserNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// build the response
		res := api_models.ReadUserResponse{
			Username: sres.Username,
			Email:    sres.Email,
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

//	@Summary		DeleteUser
//	@Description	Delete user account
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"Bearer token"
//	@Router			/me [delete]
//	@Security		Bearer
//	@Tags			User
func DeleteUser(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the user ID from the request context
		userID, ok := r.Context().Value(userIDKey).(uint)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// build the service request
		sreq := service_models.DeleteUserRequest{
			UserID: userID,
		}

		// call the service to delete user
		if err := userService.DeleteUser(sreq); err != nil {
			if errors.Is(err, services.ErrUserNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

//	@Summary		UpdateUser
//	@Description	Update user info
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string							true	"Bearer token"
//	@Param			user			body	api_models.UpdateUserRequest	true	"User info"
//	@Router			/me [put]
//	@Security		Bearer
//	@Tags			User
func UpdateUser(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the user ID from the request context
		userID, ok := r.Context().Value(userIDKey).(uint)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// decode the request body
		var req api_models.UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// build the service request
		sreq := service_models.UpdateUserRequest{
			UserID:   userID,
			Username: req.Username,
			Email:    req.Email,
		}

		// call the service to update user info
		if err := userService.UpdateUser(sreq); err != nil {
			if errors.Is(err, services.ErrUserNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
