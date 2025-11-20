package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sergey/work-track-backend/internal/models"
	"github.com/sergey/work-track-backend/internal/repository"
	"github.com/sergey/work-track-backend/internal/util"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUnauthorized       = errors.New("unauthorized access")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *models.UserRegistration) (*models.AuthResponse, error) {
	// Validate input
	if req.Login == "" || req.Password == "" {
		return nil, errors.New("login and password are required")
	}

	if req.FirstName == "" || req.LastName == "" {
		return nil, errors.New("first name and last name are required")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Avatar:       req.Avatar,
		Login:        req.Login,
		PasswordHash: hashedPassword,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := util.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, req *models.UserLogin) (*models.AuthResponse, error) {
	// Validate input
	if req.Login == "" || req.Password == "" {
		return nil, errors.New("login and password are required")
	}

	// Find user by login
	user, err := s.userRepo.FindByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Verify password
	if err := util.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := util.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}
