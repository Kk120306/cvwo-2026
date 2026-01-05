package services

import (
	"errors"
	"os"
	"time"

	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// AuthService handles authentication business logic
type AuthService struct{}

// NewAuthService creates a new instance of AuthService
func NewAuthService() *AuthService {
	return &AuthService{}
}

// AuthInput represents the data needed for signup
type AuthInput struct {
	Username string `json:"username" binding:"required"`
}


// UserResponse represents the sanitized user data returned to clients
type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarUrl"`
	IsAdmin   bool   `json:"isAdmin"`
}

// CreateUser creates a new user in the database
func (s *AuthService) CreateUser(input AuthInput) (*models.User, error) {
	// Creating the User
	// https://gorm.io/docs/create.html
	user := models.User{
		Username: input.Username,
	}
	result := database.DB.Create(&user)

	// if there was an error during creation
	if result.Error != nil {
		return nil, errors.New("failed to create user")
	}

	return &user, nil
}

// FindUserByUsername finds a user by their username
func (s *AuthService) FindUserByUsername(username string) (*models.User, error) {
	// https://gorm.io/docs/query.html - refer to inline conditions
	// Find the user with the username
	var user models.User
	dbResult := database.DB.First(&user, "Username = ?", username)

	// https://gorm.io/docs/error_handling.html
	if dbResult.Error != nil {
		if dbResult.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("could not find user with that username")
		}
		return nil, errors.New("database error")
	}

	return &user, nil
}

// FindUserByID finds a user by their ID
func (s *AuthService) FindUserByID(userID uint) (*models.User, error) {
	var user models.User
	result := database.DB.First(&user, userID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error")
	}

	return &user, nil
}

// GenerateToken generates a JWT token for the user
func (s *AuthService) GenerateToken(userID string) (string, error) {
	// generating the JWT token
	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-New-Hmac - Read this for understanding
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 day expiration
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", errors.New("failed to create token")
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	// Parse takes the token string and a function for looking up the key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return 0, errors.New("invalid token")
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get user ID from token
		userID, ok := claims["sub"].(float64)
		if !ok {
			return 0, errors.New("invalid token claims")
		}
		return uint(userID), nil
	}

	return 0, errors.New("invalid token")
}

// ToUserResponse converts a User model to a sanitized UserResponse
func (s *AuthService) ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		IsAdmin:   user.IsAdmin,
	}
}
