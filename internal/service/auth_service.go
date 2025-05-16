package service

import (
	"errors"
	"learn/internal/config"
	"learn/internal/models"
	"learn/internal/repository"
	"learn/internal/utils"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
)

// AuthService handles authentication-related business logic
type AuthService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	config    *config.Config
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		config:    config,
	}
}

// Register registers a new user
func (s *AuthService) Register(reg *models.UserRegistration) (*models.User, error) {
	// Check if user already exists
	_, err := s.userRepo.GetByUsername(reg.Username)
	if err == nil {
		return nil, ErrUserExists
	}
	if err != repository.ErrUserNotFound {
		return nil, err
	}

	// Check if email already exists
	_, err = s.userRepo.GetByEmail(reg.Email)
	if err == nil {
		return nil, ErrUserExists
	}
	if err != repository.ErrUserNotFound {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(reg.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	now := time.Now()
	user := &models.User{
		ID:           uuid.New().String(),
		Username:     reg.Username,
		Email:        reg.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Save user
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(creds *models.UserCredentials) (*models.TokenPair, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(creds.Username)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password
	if !utils.CheckPasswordHash(creds.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	tokenPair, err := utils.GenerateTokenPair(
		user.ID,
		user.Username,
		s.config.JWT.Secret,
		s.config.JWT.AccessTokenTTL,
		s.config.JWT.RefreshTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	expiresAt := time.Now().Add(s.config.JWT.RefreshTokenTTL)
	err = s.tokenRepo.Store(user.ID, tokenPair.RefreshToken, expiresAt)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*models.TokenPair, error) {
	// Validate refresh token
	claims, err := utils.ValidateJWT(refreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Check if token exists in repository
	userID, err := s.tokenRepo.GetUserIDByToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Verify user ID from token matches user ID from repository
	if claims.UserID != userID {
		return nil, ErrInvalidToken
	}

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Delete old refresh token
	err = s.tokenRepo.DeleteByToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	tokenPair, err := utils.GenerateTokenPair(
		user.ID,
		user.Username,
		s.config.JWT.Secret,
		s.config.JWT.AccessTokenTTL,
		s.config.JWT.RefreshTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	// Store new refresh token
	expiresAt := time.Now().Add(s.config.JWT.RefreshTokenTTL)
	err = s.tokenRepo.Store(user.ID, tokenPair.RefreshToken, expiresAt)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

// Logout invalidates a refresh token
func (s *AuthService) Logout(refreshToken string) error {
	return s.tokenRepo.DeleteByToken(refreshToken)
}

// LogoutAll invalidates all refresh tokens for a user
func (s *AuthService) LogoutAll(userID string) error {
	return s.tokenRepo.DeleteAllForUser(userID)
}
