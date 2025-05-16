package utils

import (
	"learn/internal/models"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	// Test parameters
	userID := "user123"
	username := "testuser"
	secret := "test-secret-key"
	expiresIn := time.Hour

	// Generate token
	token, expirationTime, err := GenerateJWT(userID, username, secret, expiresIn)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// Check expiration time is approximately correct (within a second)
	expectedExpiration := time.Now().Add(expiresIn)
	assert.WithinDuration(t, expectedExpiration, expirationTime, time.Second)

	// Parse and verify token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	// Verify claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, userID, claims["user_id"])
	assert.Equal(t, username, claims["username"])
	
	// Verify expiration time in claims
	expClaim, ok := claims["exp"].(float64)
	assert.True(t, ok)
	assert.InDelta(t, expirationTime.Unix(), int64(expClaim), 1)
}

func TestValidateJWT(t *testing.T) {
	// Test parameters
	userID := "user123"
	username := "testuser"
	secret := "test-secret-key"
	expiresIn := time.Hour

	// Generate a valid token
	validToken, _, err := GenerateJWT(userID, username, secret, expiresIn)
	assert.NoError(t, err)

	// Generate an expired token
	expiredToken, _, err := GenerateJWT(userID, username, secret, -time.Hour)
	assert.NoError(t, err)

	// Create an invalid token (wrong signature)
	invalidToken, _, err := GenerateJWT(userID, username, "wrong-secret", expiresIn)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		tokenString string
		secret      string
		wantClaims  *models.TokenClaims
		wantErr     error
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			secret:      secret,
			wantClaims: &models.TokenClaims{
				UserID:   userID,
				Username: username,
			},
			wantErr: nil,
		},
		{
			name:        "Expired token",
			tokenString: expiredToken,
			secret:      secret,
			wantClaims:  nil,
			wantErr:     ErrExpiredToken,
		},
		{
			name:        "Invalid token (wrong signature)",
			tokenString: invalidToken,
			secret:      secret,
			wantClaims:  nil,
			wantErr:     ErrInvalidToken,
		},
		{
			name:        "Empty token",
			tokenString: "",
			secret:      secret,
			wantClaims:  nil,
			wantErr:     jwt.ErrTokenMalformed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateJWT(tt.tokenString, tt.secret)
			
			if tt.wantErr != nil {
				assert.Error(t, err)
				if tt.wantErr == ErrExpiredToken {
					assert.ErrorIs(t, err, ErrExpiredToken)
				} else if tt.wantErr == ErrInvalidToken {
					assert.Error(t, err) // The error might not be exactly ErrInvalidToken
				}
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.wantClaims.UserID, claims.UserID)
				assert.Equal(t, tt.wantClaims.Username, claims.Username)
			}
		})
	}
}

func TestGenerateTokenPair(t *testing.T) {
	// Test parameters
	userID := "user123"
	username := "testuser"
	secret := "test-secret-key"
	accessTTL := time.Hour
	refreshTTL := time.Hour * 24 * 7 // 1 week

	// Generate token pair
	tokenPair, err := GenerateTokenPair(userID, username, secret, accessTTL, refreshTTL)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	
	// Verify expiration time is set correctly
	expectedExpiration := time.Now().Add(accessTTL).Unix()
	assert.InDelta(t, expectedExpiration, tokenPair.ExpiresAt, 1)

	// Validate access token
	accessClaims, err := ValidateJWT(tokenPair.AccessToken, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, accessClaims.UserID)
	assert.Equal(t, username, accessClaims.Username)

	// Validate refresh token
	refreshClaims, err := ValidateJWT(tokenPair.RefreshToken, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, refreshClaims.UserID)
	assert.Equal(t, username, refreshClaims.Username)
}

func TestInvalidSigningMethod(t *testing.T) {
	// Create a token with a different signing method
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"user_id":  "user123",
		"username": "testuser",
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)
	
	// Try to validate the token
	claims, err := ValidateJWT(tokenString, "test-secret")
	
	// Should fail with invalid token error
	assert.Error(t, err)
	assert.Nil(t, claims)
}