package middleware

import (
	"net/http"

	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/gin-gonic/gin"
)

// Middleware to check if the user is an admin.
// Must be run in subsequent to CheckAuth middleware.
func CheckAdmin(c *gin.Context) {
	user, exist := c.Get("user") // Can assume user always exists because CheckAuth is ran before
	if !exist {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Map the user to models
	u := user.(models.User)
	// if the "IsAdmin" is false we send a forbidden stat
	if !u.IsAdmin {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}
