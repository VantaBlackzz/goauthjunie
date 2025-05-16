package handlers

import (
	"learn/internal/middleware"
	"learn/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userRepo repository.UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// GetProfile handles getting the user's profile
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User profile"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(router *gin.Engine, middleware gin.HandlerFunc) {
	user := router.Group("/user")
	user.Use(middleware)
	{
		user.GET("/profile", h.GetProfile)
	}
}
