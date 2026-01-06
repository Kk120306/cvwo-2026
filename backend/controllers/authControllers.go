package controllers

import (
	"net/http"

	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/Kk120306/cvwo-2026/backend/services"
	"github.com/gin-gonic/gin"
)

// AuthController handles HTTP requests for authentication
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

// Signup function - handles user registration
func (ac *AuthController) Signup(c *gin.Context) {
	// structure of the request body that we need
	var body struct {
		Username string `json:"username" binding:"required"`
	}

	// if parsing the req to body fails, returns non nil, sends bad request
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body, please provide valid details",
		})
		return
	}

	// Create user through service layer
	user, err := ac.authService.CreateUser(services.AuthInput{
		Username: body.Username,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Successfully created a user so we just return a response here
	c.JSON(http.StatusOK, gin.H{
		"user": ac.authService.ToUserResponse(user),
	})
}

// Login function - handles user authentication
func (ac *AuthController) Login(c *gin.Context) {
	// structure the req body
	var body struct {
		Username string `json:"username" binding:"required"`
	}

	// compare the parsed input with the structure
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body. Please provide valid details",
		})
		return
	}

	// Find user through service layer
	user, err := ac.authService.FindUserByUsername(body.Username)
	if err != nil {
		// Determine status code based on error type
		statusCode := http.StatusInternalServerError
		if err.Error() == "could not find user with that username" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Generate token through service layer
	tokenString, err := ac.authService.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Creating a cookie with the token
	// https://krisnacahyono.medium.com/api-authentication-with-go-481f87947c26
	c.SetSameSite(http.SameSiteLaxMode) // makes sure cookies are not sent on cross-site requests
	c.SetCookie(
		"Authorization", // Cookie name
		tokenString,     // Cookies JWT value
		3600*24*30,      // age
		"",              // path - accessible everywhere
		"",              // domain - only sent to host
		false,           // secure on dev so it works on localhost TODO: set true later
		true,            // Https only
	)

	// Sending success response without password for security reasons
	c.JSON(http.StatusOK, gin.H{
		"user": ac.authService.ToUserResponse(user),
	})
}

// Validate function - sends user data and checks if token is valid
func (ac *AuthController) Validate(c *gin.Context) {
	// Retrieve user from middleware
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// checks to ensure that user type is models.User so it can be refactored for response
	user, ok := u.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}

	// Return a sanitized user object without password
	c.JSON(http.StatusOK, gin.H{
		"user": ac.authService.ToUserResponse(&user),
	})
}

// Logout Function - clears the cookie
func (ac *AuthController) Logout(c *gin.Context) {
	c.SetCookie(
		"Authorization", // Cookie name
		"",              // Cookies JWT value
		-1,              // age - ensures that the cookie expires immediately
		"",              // path - accessible everywhere
		"",              // domain - only sent to host
		false,           // secure on dev so it works on localhost TODO: set true later
		true,            // Https only
	)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
