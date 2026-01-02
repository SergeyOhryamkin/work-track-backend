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
	ErrSessionNotFound    = errors.New("session not found")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.UserSessionRepository
	jwtSecret   string
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.UserSessionRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
	}
}

// Register creates a new user account and starts a session
func (s *AuthService) Register(ctx context.Context, req *models.UserRegistration, meta *models.SessionMetadata) (*models.AuthResponse, error) {
	if req.Login == "" || req.Password == "" {
		return nil, errors.New("login and password are required")
	}
	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Login:        req.Login,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := util.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	sessionID, err := s.createSession(ctx, user.ID, refreshToken, meta)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &models.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
		SessionID:    sessionID,
	}, nil
}

// Login authenticates a user and starts a session
func (s *AuthService) Login(ctx context.Context, req *models.UserLogin, meta *models.SessionMetadata) (*models.AuthResponse, error) {
	if req.Login == "" || req.Password == "" {
		return nil, errors.New("login and password are required")
	}

	user, err := s.userRepo.FindByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if err := util.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := util.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	sessionID, err := s.createSession(ctx, user.ID, refreshToken, meta)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &models.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
		SessionID:    sessionID,
	}, nil
}

// Logout completes an active session
func (s *AuthService) Logout(ctx context.Context, userID, sessionID int) error {
	if sessionID == 0 {
		return errors.New("session_id is required")
	}

	_, err := s.sessionRepo.Complete(ctx, sessionID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return ErrSessionNotFound
		}
		return fmt.Errorf("failed to complete session: %w", err)
	}

	return nil
}

// RefreshToken generates a new access token and refresh token using a valid refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	// Validate refresh token
	claims, err := util.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, ErrUnauthorized
	}

	// Find active session with this refresh token
	session, err := s.sessionRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to find session: %w", err)
	}

	// Extra security check: ensure token belongs to the user in the session
	if claims.UserID != session.UserID {
		return nil, ErrUnauthorized
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Generate new tokens
	newToken, err := util.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	newRefreshToken, err := util.GenerateRefreshToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update session with new refresh token
	err = s.sessionRepo.UpdateRefreshToken(ctx, session.ID, newRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to update session refresh token: %w", err)
	}

	return &models.AuthResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
		User:         *user,
		SessionID:    session.ID,
	}, nil
}

func (s *AuthService) createSession(ctx context.Context, userID int, refreshToken string, meta *models.SessionMetadata) (int, error) {
	if meta == nil {
		meta = &models.SessionMetadata{}
	}

	session := &models.UserSession{
		UserID:       userID,
		RefreshToken: refreshToken,
		Device:       meta.Device,
		Platform:     meta.Platform,
		UserAgent:    meta.UserAgent,
		IPAddress:    meta.IPAddress,
	}

	id, err := s.sessionRepo.Create(ctx, session)
	if err != nil {
		return 0, err
	}

	return id, nil
}
