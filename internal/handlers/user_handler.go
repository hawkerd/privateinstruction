package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hawkerd/privateinstruction/internal/models/api_models"
	"github.com/hawkerd/privateinstruction/internal/services"
)

const userIDKey = "userID"

// get user info
func ReadUser(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the user ID from the request context
		userID, ok := r.Context().Value(userIDKey).(uint)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// call the service to get user info
		user, err := userService.ReadUser(userID)
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
			Username: user.Username,
			Email:    user.Email,
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
