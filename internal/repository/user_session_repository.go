package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sergey/work-track-backend/internal/models"
)

var (
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionAlreadyClosed = errors.New("session already closed")
)

// UserSessionRepository handles persistence for user session records
type UserSessionRepository struct {
	db *sql.DB
}

// NewUserSessionRepository creates a new repository instance
func NewUserSessionRepository(db *sql.DB) *UserSessionRepository {
	return &UserSessionRepository{db: db}
}

// Create inserts a new session record and returns the session ID
func (r *UserSessionRepository) Create(ctx context.Context, session *models.UserSession) (int, error) {
	query := `
		INSERT INTO user_sessions (user_id, device, platform, user_agent, ip_address, refresh_token)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		session.UserID,
		session.Device,
		session.Platform,
		session.UserAgent,
		session.IPAddress,
		session.RefreshToken,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get session id: %w", err)
	}

	session.ID = int(id)

	// Load timestamps
	err = r.db.QueryRowContext(ctx, `
		SELECT login_at, created_at
		FROM user_sessions
		WHERE id = ?
	`, session.ID).Scan(&session.LoginAt, &session.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch session timestamps: %w", err)
	}

	return session.ID, nil
}

// FindByRefreshToken retrieves a session by its refresh token
func (r *UserSessionRepository) FindByRefreshToken(ctx context.Context, token string) (*models.UserSession, error) {
	query := `
		SELECT id, user_id, login_at, logout_at, refresh_token, device, platform, user_agent,
		       ip_address, session_duration_seconds, created_at
		FROM user_sessions
		WHERE refresh_token = ? AND logout_at IS NULL
	`

	var session models.UserSession
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.LoginAt,
		&session.LogoutAt,
		&session.RefreshToken,
		&session.Device,
		&session.Platform,
		&session.UserAgent,
		&session.IPAddress,
		&session.SessionDurationSeconds,
		&session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to find session by refresh token: %w", err)
	}

	return &session, nil
}

// UpdateRefreshToken updates the refresh token for a session
func (r *UserSessionRepository) UpdateRefreshToken(ctx context.Context, sessionID int, token string) error {
	query := `UPDATE user_sessions SET refresh_token = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, token, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update refresh token: %w", err)
	}
	return nil
}
func (r *UserSessionRepository) Complete(ctx context.Context, sessionID, userID int) (*models.UserSession, error) {
	query := `
		UPDATE user_sessions
		SET
			logout_at = CURRENT_TIMESTAMP,
			session_duration_seconds = CAST(
				(strftime('%s', CURRENT_TIMESTAMP) - strftime('%s', login_at)) AS INTEGER
			)
		WHERE id = ? AND user_id = ? AND logout_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, sessionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to complete session: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch affected rows: %w", err)
	}
	if rows == 0 {
		return nil, ErrSessionNotFound
	}

	var session models.UserSession
	err = r.db.QueryRowContext(ctx, `
		SELECT id, user_id, login_at, logout_at, device, platform, user_agent,
		       ip_address, session_duration_seconds, created_at
		FROM user_sessions
		WHERE id = ?
	`, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.LoginAt,
		&session.LogoutAt,
		&session.Device,
		&session.Platform,
		&session.UserAgent,
		&session.IPAddress,
		&session.SessionDurationSeconds,
		&session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to load completed session: %w", err)
	}

	return &session, nil
}
