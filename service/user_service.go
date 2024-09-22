package service

import (
	"context"
	"errors"
	"github.com/trodix/go-rest-api/models"
	"github.com/trodix/go-rest-api/repository"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// CreateUser creates a new user after applying business logic (e.g., checks, transformations)
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	// Example of additional business logic
	if len(user.Username) < 4 {
		return errors.New("username must be at least 4 characters")
	}
	return s.Repo.CreateUser(ctx, user)
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.Repo.GetUserByID(ctx, id)
}

// GetUserByID retrieves a user by their Username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.Repo.GetUserByUsername(ctx, username)
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.Repo.GetAllUsers(ctx)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.Repo.UpdateUser(ctx, user)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.Repo.DeleteUser(ctx, id)
}
