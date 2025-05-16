package middleware

import (
	"learn/internal/config"
	"learn/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware creates a middleware for JWT authentication
func JWTMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := parts[1]

		// Validate the token
		claims, err := utils.ValidateJWT(tokenString, config.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the user ID and username in the context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// GetUserID gets the user ID from the context
func GetUserID(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	return userID.(string)
}

// GetUsername gets the username from the context
func GetUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}