package controllers

import (
	"net/http"

	"github.com/Kk120306/cvwo-2026/backend/services"
	"github.com/gin-gonic/gin"
)

// UserController handles HTTP requests for users
type UserController struct {
	userService *services.UserService
}

// NewUserController creates a new instance of UserController
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// GetUserProfile retrieves a user's profile with optional posts and comments
func (uc *UserController) GetUserProfile(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username required"})
		return
	}

	// Check query parameters for optional data
	includePosts := c.Query("posts") == "true"
	includeComments := c.Query("comments") == "true"

	// Get user profile through service layer
	profile, err := uc.userService.GetUserProfile(username, includePosts, includeComments)
	if err != nil {
		// Determine status code based on error type
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": profile})
}
