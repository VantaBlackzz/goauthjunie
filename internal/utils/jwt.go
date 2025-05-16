package utils

import (
	"errors"
	"learn/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// GenerateJWT generates a new JWT token
func GenerateJWT(userID, username, secret string, expiresIn time.Duration) (string, time.Time, error) {
	expirationTime := time.Now().Add(expiresIn)
	
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	
	return tokenString, expirationTime, nil
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString, secret string) (*models.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, err
	}
	
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}
	
	username, ok := claims["username"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}
	
	return &models.TokenClaims{
		UserID:   userID,
		Username: username,
	}, nil
}

// GenerateTokenPair generates both access and refresh tokens
func GenerateTokenPair(userID, username, secret string, accessTTL, refreshTTL time.Duration) (*models.TokenPair, error) {
	// Generate access token
	accessToken, expiresAt, err := GenerateJWT(userID, username, secret, accessTTL)
	if err != nil {
		return nil, err
	}
	
	// Generate refresh token
	refreshToken, _, err := GenerateJWT(userID, username, secret, refreshTTL)
	if err != nil {
		return nil, err
	}
	
	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}