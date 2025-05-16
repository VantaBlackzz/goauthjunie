package models

import (
	"time"
)

// TokenPair represents a pair of JWT tokens (access and refresh)
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"` // Unix timestamp for access token expiration
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

// RefreshRequest represents a request to refresh an access token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// TokenResponse represents the response for token operations
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}