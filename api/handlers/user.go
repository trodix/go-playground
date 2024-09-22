package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/trodix/go-rest-api/api/middleware"
	"github.com/trodix/go-rest-api/models"
	"github.com/trodix/go-rest-api/service"
	"github.com/trodix/go-rest-api/utils"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.(*middleware.ResponseRecorder).CaptureError(http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate user
	if err := utils.ValidateStruct(user); err != nil {
		w.(*middleware.ResponseRecorder).CaptureError(http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	err := h.Service.CreateUser(r.Context(), &user)
	if err != nil {
		w.(*middleware.ResponseRecorder).CaptureError(http.StatusBadRequest, "Failed to create user: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.(*middleware.ResponseRecorder).CaptureError(http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.Service.GetUserByID(r.Context(), id)
	if err != nil {
		w.(*middleware.ResponseRecorder).CaptureError(http.StatusNotFound, "User not found")
		return
	}

	json.NewEncoder(w).Encode(user)
}

// GetUsers retrieves all users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers(r.Context())
	if err != nil {
		w.(*middleware.ResponseRecorder).CaptureError(http.StatusInternalServerError, "Failed to retrieve users")
		return
	}
	if users == nil {
		users = []*models.User{}
	}
	json.NewEncoder(w).Encode(users)
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.(*middleware.ResponseRecorder).CaptureError(http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.(*middleware.ResponseRecorder).CaptureError(http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate user
	if err := utils.ValidateStruct(user); err != nil {
		w.(*middleware.ResponseRecorder).CaptureError(http.StatusBadRequest, err.Error())
		return
	}

	user.ID = id
	if err := h.Service.UpdateUser(r.Context(), &user); err != nil {
		w.(*middleware.ResponseRecorder).CaptureError(http.StatusInternalServerError, "Failed to update user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser deletes a user by ID
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.(*middleware.ResponseRecorder).CaptureError(http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.Service.DeleteUser(r.Context(), id); err != nil {
		w.(*middleware.ResponseRecorder).CaptureError(http.StatusInternalServerError, "Failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
