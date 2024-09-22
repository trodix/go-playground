package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/trodix/go-rest-api/models"
	"github.com/trodix/go-rest-api/repository"
	"github.com/trodix/go-rest-api/utils"
)

type UserHandler struct {
    Repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
    return &UserHandler{Repo: repo}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validate user
    if err := utils.ValidateStruct(user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := h.Repo.CreateUser(r.Context(), &user)
    if err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
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
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    user, err := h.Repo.GetUserByID(r.Context(), id)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(user)
}

// GetUsers retrieves all users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.Repo.GetAllUsers(r.Context())
    if err != nil {
		log.Printf("Failed to retrieve users: %s", err.Error())
        http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
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
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validate user
    if err := utils.ValidateStruct(user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user.ID = id // Set the ID for the update
    if err := h.Repo.UpdateUser(r.Context(), &user); err != nil {
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// DeleteUser deletes a user by ID
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    if err := h.Repo.DeleteUser(r.Context(), id); err != nil {
        http.Error(w, "Failed to delete user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent) // 204 No Content
}
