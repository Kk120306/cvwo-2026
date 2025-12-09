package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
			"error": "Failed to read request body, please provide valid details",
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
	result := database.DB.Create(&user)

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

// Login function 
func Login(c *gin.Context) {
	// structure the req body
	var body struct {
		Email    string
		Password string
	}

	// compare the parsed input with the structure
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body. Please provide valid details",
		})
		return
	}

	// https://gorm.io/docs/query.html - refer to inline conditions
	// Find the user with the email
	var user models.User
	database.DB.First(&user, "Email = ?", body.Email)
	// SELECT * FROM users WHERE id = 'string_primary_key';
	if user.ID == "0" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Comparing password using bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// generating the JWT token
	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-New-Hmac - Read this for understanding
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 day expiration
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token. Try again.",
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

	// Sending success response
	c.JSON(http.StatusOK, gin.H{})

}
