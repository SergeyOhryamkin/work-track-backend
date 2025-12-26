package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/sergey/work-track-backend/internal/models"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (first_name, last_name, avatar, login, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`

	result, err := r.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.Avatar, user.Login, user.PasswordHash)
	if err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	user.ID = int(id)

	// Fetch created_at and updated_at
	err = r.db.QueryRowContext(ctx, "SELECT created_at, updated_at FROM users WHERE id = ?", user.ID).
		Scan(&user.CreatedAt, &user.UpdatedAt)

	return nil
}

// FindByLogin retrieves a user by login
func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, avatar, login, password_hash, created_at, updated_at
		FROM users
		WHERE login = ?
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, login).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Avatar, &user.Login, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by login: %w", err)
	}

	return &user, nil
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, avatar, login, password_hash, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Avatar, &user.Login, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return &user, nil
}
