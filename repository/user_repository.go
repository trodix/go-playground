package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/trodix/go-rest-api/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRow(ctx, query, user.Username, user.Email).Scan(&user.ID)
	return err
}

// GetUserByID retrieves a user by ID from the database
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := "SELECT id, username, email FROM users WHERE id = $1"
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by Username from the database
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := "SELECT id, username, email FROM users WHERE username = $1"
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAllUsers retrieves all users from the database
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	query := "SELECT id, username, email FROM users"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// UpdateUser updates an existing user in the database
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := "UPDATE users SET username = $1, email = $2 WHERE id = $3"
	_, err := r.db.Exec(ctx, query, user.Username, user.Email, user.ID)
	return err
}

// DeleteUser deletes a user by ID from the database
func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(ctx, query, id)
	return err
}
