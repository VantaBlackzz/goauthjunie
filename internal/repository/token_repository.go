package repository

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrTokenNotFound      = errors.New("token not found")
	ErrTokenAlreadyExists = errors.New("token already exists")
	ErrTokenExpired       = errors.New("token expired")
)

// TokenRepository defines the interface for token data access
type TokenRepository interface {
	Store(userID, refreshToken string, expiresAt time.Time) error
	GetUserIDByToken(refreshToken string) (string, error)
	DeleteByToken(refreshToken string) error
	DeleteAllForUser(userID string) error
}

// InMemoryTokenRepository implements TokenRepository with an in-memory store
type InMemoryTokenRepository struct {
	// Map of refresh tokens to user IDs and expiration times
	tokens map[string]tokenInfo
	mutex  sync.RWMutex
}

type tokenInfo struct {
	UserID    string
	ExpiresAt time.Time
}

// NewInMemoryTokenRepository creates a new in-memory token repository
func NewInMemoryTokenRepository() *InMemoryTokenRepository {
	return &InMemoryTokenRepository{
		tokens: make(map[string]tokenInfo),
	}
}

// Store adds a new refresh token to the repository
func (r *InMemoryTokenRepository) Store(userID, refreshToken string, expiresAt time.Time) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tokens[refreshToken]; exists {
		return ErrTokenAlreadyExists
	}

	r.tokens[refreshToken] = tokenInfo{
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	return nil
}

// GetUserIDByToken retrieves a user ID by refresh token
func (r *InMemoryTokenRepository) GetUserIDByToken(refreshToken string) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	info, exists := r.tokens[refreshToken]
	if !exists {
		return "", ErrTokenNotFound
	}

	// Check if token is expired
	if time.Now().After(info.ExpiresAt) {
		return "", ErrTokenExpired
	}

	return info.UserID, nil
}

// DeleteByToken removes a refresh token
func (r *InMemoryTokenRepository) DeleteByToken(refreshToken string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tokens[refreshToken]; !exists {
		return ErrTokenNotFound
	}

	delete(r.tokens, refreshToken)
	return nil
}

// DeleteAllForUser removes all refresh tokens for a user
func (r *InMemoryTokenRepository) DeleteAllForUser(userID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for token, info := range r.tokens {
		if info.UserID == userID {
			delete(r.tokens, token)
		}
	}
	return nil
}