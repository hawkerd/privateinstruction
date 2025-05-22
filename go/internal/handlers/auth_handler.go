package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hawkerd/privateinstruction/internal/auth"
	"github.com/hawkerd/privateinstruction/internal/models/api_models"
	"github.com/hawkerd/privateinstruction/internal/models/service_models"
	"github.com/hawkerd/privateinstruction/internal/services"
)

// @Summary		Sign Up
// @Description	Sign up a new user
// @Accept			json
// @Produce		json
// @Param			user	body		api_models.SignUpRequest	true	"User details for sign up"
// @Success		201		{string}	string						"User created successfully"
// @Failure		400		{string}	string						"Bad Request"
// @Failure		409		{string}	string						"Conflict"
// @Failure		500		{string}	string						"Internal Server Error"
// @Router			/signup [post]
// @Tags			Auth
func SignUp(authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decode the request body
		var req api_models.SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Please fill out all fields", http.StatusBadRequest)
			return
		}

		// input validation
		if req.Username == "" || req.Password == "" || req.Email == "" {
			http.Error(w, "Please fill out all fields", http.StatusBadRequest)
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
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
			}
			return
		}

		// encode the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}

//		@Summary		Sign In
//		@Description	Sign in an existing user with username/email and password
//	 @Description	Also sets the refresh token in the cookie
//		@Accept			json
//		@Produce		json
//		@Param			user	body	api_models.SignInRequest	true	"User credentials for sign in"
//		@Router			/signin [post]
//		@Tags			Auth
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

		// set the refresh token in the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    sres.RefreshToken,
			Expires:  sres.RefreshTokenExpiration,
			HttpOnly: true,
			Secure:   false, // for testing
			Path:     "/auth/refresh",
			SameSite: http.SameSiteStrictMode,
		})

		// build the response
		res := api_models.SignInResponse{
			AccessToken: sres.AccessToken,
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

// @Summary		Update Password
// @Description	Update the password for an existing user
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			Authorization	header	string								true	"Bearer Token"
// @Param			user			body	api_models.UpdatePasswordRequest	true	"User credentials for updating password"
// @Router			/me/password [put]
// @Tags			Auth
func UpdatePassword(authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the user ID from the request context
		userID, ok := r.Context().Value(userIDKey).(uint)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// decode the request body
		var req api_models.UpdatePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// input validation
		if userID == 0 || req.OldPassword == "" || req.NewPassword == "" {
			http.Error(w, "user_id, old_password, and new_password are required", http.StatusBadRequest)
			return
		}

		// build the service request
		sreq := service_models.UpdatePasswordRequest{
			UserID:      userID,
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

// @Summary		Refresh Access Token
// @Description	Refresh the access token using the refresh token
// @Accept		json
// @Produce		json
// @Param			Authorization	header	string	true	"Bearer Token"
// @Router			/auth/refresh [post]
// @Tags			Auth
func RefreshToken(authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the JWT from the request header
		tokenString, err := auth.ExtractJWT(r)
		if err != nil {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		// parse the token to extract the user ID
		userID, err := auth.ParseID(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// extract the refresh token from the cookie
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		refreshToken := cookie.Value

		// input validation
		if refreshToken == "" {
			http.Error(w, "token is required", http.StatusBadRequest)
			return
		}

		// build the service request
		sreq := service_models.RefreshTokenRequest{
			RefreshToken: refreshToken,
			UserID:       userID,
		}

		// call the service
		sres, err := authService.RefreshAccessToken(sreq)
		if err != nil {
			if errors.Is(err, services.ErrInvalidCredentials) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		// set the refresh token in the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    sres.RefreshToken,
			Expires:  sres.RefreshTokenExpiration,
			HttpOnly: true,
			Secure:   false, // for testing
			Path:     "/auth/refresh",
			SameSite: http.SameSiteStrictMode,
		})

		// build the response
		res := api_models.RefreshTokenResponse{
			AccessToken: sres.AccessToken,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}
