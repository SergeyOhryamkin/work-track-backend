package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Avatar       string    `json:"avatar,omitempty"` // URL or path to avatar image
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserRegistration represents the data needed to register a new user
type UserRegistration struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// UserLogin represents the data needed to log in
type UserLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
