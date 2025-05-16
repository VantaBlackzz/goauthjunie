package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	// In a real application, you would include database connection details
	// For this example, we'll use an in-memory store
	InMemory bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Server config
	port := getEnv("SERVER_PORT", "8080")
	readTimeout, _ := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "10"))
	writeTimeout, _ := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "10"))

	// JWT config
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	accessTokenTTL, _ := strconv.Atoi(getEnv("JWT_ACCESS_TOKEN_TTL", "15"))      // 15 minutes
	refreshTokenTTL, _ := strconv.Atoi(getEnv("JWT_REFRESH_TOKEN_TTL", "10080")) // 7 days

	config := &Config{
		Server: ServerConfig{
			Port:         port,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
		},
		JWT: JWTConfig{
			Secret:          jwtSecret,
			AccessTokenTTL:  time.Duration(accessTokenTTL) * time.Minute,
			RefreshTokenTTL: time.Duration(refreshTokenTTL) * time.Minute,
		},
		Database: DatabaseConfig{
			InMemory: true,
		},
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
