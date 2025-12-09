package controllers

import (
	"net/http"

	"github.com/Kk120306/cvwo-2026/backend/helpers"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Sign up function
func Signup(c *gin.Context) {
	// structure of the request body that we need
	var body struct {
		Username string
		Email    string
		Password string
	}

	// if parsing the req to body fails, returns non nil, sends bad request
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// Using Bcrypt to hash
	// checking if there was an error during encryption
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to encrypt password",
		})
		return
	}

	// Creating the User
	// https://gorm.io/docs/create.html
	user := models.User{
		Username:          body.Username,
		Email:             body.Email,
		EncryptedPassword: string(encryptedPass),
	}
	result := helpers.DB.Create(&user)

	// if there was an error during creation
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Successfully created a user so we just return a response here
	c.JSON(http.StatusOK, gin.H{})
}
