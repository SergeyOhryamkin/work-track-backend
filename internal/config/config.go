package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
	S3       S3Config
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Path string // SQLite database file path
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret string
}

// CORSConfig holds CORS-related configuration
type CORSConfig struct {
	AllowedOrigins []string
}

// S3Config holds S3-related configuration
type S3Config struct {
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	Region    string
	Prefix    string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	allowedOrigins := strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Path: getEnv("DB_PATH", "./worktrack.db"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
		CORS: CORSConfig{
			AllowedOrigins: allowedOrigins,
		},
		S3: S3Config{
			Endpoint:  getEnv("S3_ENDPOINT", ""),
			Bucket:    getEnv("S3_BUCKET", ""),
			AccessKey: getEnv("S3_ACCESS_KEY", ""),
			SecretKey: getEnv("S3_SECRET_KEY", ""),
			Region:    getEnv("S3_REGION", ""),
			Prefix:    getEnv("S3_PREFIX", ""),
		},
	}

	// Validate required fields
	if config.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return config, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ConnectionString returns the SQLite connection string
func (c *DatabaseConfig) ConnectionString() string {
	return c.Path
}
