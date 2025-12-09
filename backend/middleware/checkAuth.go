package middleware

import (
	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"time"
)

func CheckAuth(c *gin.Context) {
	// Get the cookie in the request
	tokenString, err := c.Cookie("Authorization")
	if err != nil || tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-Parse-Hmac
	// validating the token - Taken from the documentation above
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok { //  If we have access
		// Check the age
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Now find the user from the token
		// https://gorm.io/docs/query.html - refer to retrieving by primary key
		userID := claims["sub"].(string) // extracting userID from token
		var user models.User
		// https://gorm.io/docs/error_handling.html - .Error, cant just get err as second return value
		err := database.DB.First(&user, "id = ?", userID).Error
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Attach the user to the request context
		c.Set("user", user)

		c.Next()
	} else { // If we don't have access
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Next()
}
