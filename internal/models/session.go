package models

import "time"

// UserSession represents a single login/logout activity record
type UserSession struct {
	ID                     int        `json:"id"`
	UserID                 int        `json:"user_id"`
	LoginAt                time.Time  `json:"login_at"`
	LogoutAt               *time.Time `json:"logout_at,omitempty"`
	Device                 string     `json:"device,omitempty"`
	Platform               string     `json:"platform,omitempty"`
	UserAgent              string     `json:"user_agent,omitempty"`
	IPAddress              string     `json:"ip_address,omitempty"`
	SessionDurationSeconds *int       `json:"session_duration_seconds,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
}

// SessionMetadata contains contextual data captured during login/logout
type SessionMetadata struct {
	Device    string
	Platform  string
	UserAgent string
	IPAddress string
}
