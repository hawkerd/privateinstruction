package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/hawkerd/privateinstruction/internal/models/api_models"
	"github.com/hawkerd/privateinstruction/internal/models/service_models"
	"github.com/hawkerd/privateinstruction/internal/services"
)

// helper function to extract the class ID from the request
func getClassIDFromRequest(r *http.Request) (uint, error) {
	classIDStr := chi.URLParam(r, "id")
	if classIDStr == "" {
		return 0, errors.New("class ID is required")
	}

	classID, err := strconv.ParseUint(classIDStr, 10, 32)
	if err != nil {
		return 0, errors.New("invalid class ID")
	}

	return uint(classID), nil
}

//	@Summary		CreateClass
//	@Description	Create a new class
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string							true	"Bearer token"
//	@Param			class			body	api_models.CreateClassRequest	true	"Class info"
//	@Router			/class [post]
//	@Security		Bearer
//	@Tags			Class
func CreateClass(classService *services.ClassService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the user ID from the request context
		userID, ok := r.Context().Value(userIDKey).(uint)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// decode the request body
		var req api_models.CreateClassRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// input validation
		if req.Name == "" {
			req.Name = "Unnamed Class"
		}

		// build the service request
		sreq := service_models.CreateClassRequest{
			Name:        req.Name,
			Description: req.Description,
			UserID:      userID,
		}

		// call the service
		if err := classService.CreateClass(sreq); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

//	@Summary		DeleteClass
//	@Description	Delete a class
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"Bearer token"
//	@Param			id				path	int		true	"Class ID"
//	@Router			/class/{id} [delete]
//	@Security		Bearer
//	@Tags			Class
func DeleteClass(classService *services.ClassService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the user ID from the request context
		userID, ok := r.Context().Value(userIDKey).(uint)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// extract the class ID from the URL
		classID, err := getClassIDFromRequest(r)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// build the service request
		sreq := service_models.DeleteClassRequest{
			ClassID: classID,
			UserID:  userID,
		}

		// call the service
		err = classService.DeleteClass(sreq)
		if err != nil {
			if errors.Is(err, services.ErrClassNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if errors.Is(err, services.ErrUnauthorized) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// send a 204 No Content response
		w.WriteHeader(http.StatusNoContent)
	}
}

//	@Summary		ReadClass
//	@Description	Read a class by ID
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"Bearer token"
//	@Param			id				path	int		true	"Class ID"
//	@Router			/class/{id} [get]
//	@Security		Bearer
//	@Tags			Class
func ReadClass(classService *services.ClassService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the user ID from the request context
		userID, ok := r.Context().Value(userIDKey).(uint)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// extract the class ID from the URL
		classID, err := getClassIDFromRequest(r)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// build the service request
		sreq := service_models.ReadClassRequest{
			ClassID: classID,
			UserID:  userID,
		}

		// call the service
		sres, err := classService.ReadClass(sreq)
		if err != nil {
			if errors.Is(err, services.ErrClassNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if errors.Is(err, services.ErrUnauthorized) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// build the response
		res := api_models.ReadClassResponse{
			Name:        sres.Name,
			Description: sres.Description,
			CreatedAt:   sres.CreatedAt,
			CreatedBy:   sres.CreatedBy,
		}

		// encode the response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}

//	@Summary		UpdateClass
//	@Description	Update a class
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string							true	"Bearer token"
//	@Param			id				path	int								true	"Class ID"
//	@Param			class			body	api_models.UpdateClassRequest	true	"Class info"
//	@Router			/class/{id} [put]
//	@Security		Bearer
//	@Tags			Class
func UpdateClass(classService *services.ClassService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the user ID from the request context
		userID, ok := r.Context().Value(userIDKey).(uint)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// extract the class ID from the URL
		classID, err := getClassIDFromRequest(r)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// decode the request body
		var req api_models.UpdateClassRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// input validation
		if req.Name == "" {
			req.Name = "Unnamed Class"
		}

		// build the service request
		sreq := service_models.UpdateClassRequest{
			ClassID:     classID,
			Name:        req.Name,
			Description: req.Description,
			UserID:      userID,
		}
		// call the service
		err = classService.UpdateClass(sreq)
		if err != nil {
			if errors.Is(err, services.ErrClassNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if errors.Is(err, services.ErrUnauthorized) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// send a 204 No Content response
		w.WriteHeader(http.StatusNoContent)
	}
}
