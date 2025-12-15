package middleware

import (
	"os"
	"time"

	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// OptionalAuth is middleware that checks for authentication but doesn't require it
// If a valid token is present, it sets the user in context
// If no token or invalid token, it continues without setting user
// Used for dashboard where users can be either logged in or not and still have access
// Logic is the same as checkAuth - check documentation attached there
func OptionalAuth(c *gin.Context) {
	// Get the cookie in the request
	tokenString, err := c.Cookie("Authorization")
	if err != nil || tokenString == "" {
		c.Next()
		return
	}

	// validating the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		c.Next()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok { // If we have access
		// Check the age
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.Next()
			return
		}

		// Now find the user from the token
		userID := claims["sub"].(string) // extracting userID from token
		var user models.User
		err := database.DB.First(&user, "id = ?", userID).Error
		if err != nil {
			c.Next()
			return
		}

		// Attach the user to the request context
		c.Set("user", user)

		c.Next()
	} else { // If we don't have access
		c.Next()
	}

	c.Next()
}
