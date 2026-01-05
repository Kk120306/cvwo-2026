package services

import (
	"errors"

	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"gorm.io/gorm"
)

// UserService handles user business logic
type UserService struct{}

// NewUserService creates a new instance of UserService
func NewUserService() *UserService {
	return &UserService{}
}

// UserProfile represents a user's profile with statistics
type UserProfile struct {
	models.User
	PostCount    int64            `json:"postCount"`
	CommentCount int64            `json:"commentCount"`
	Posts        []models.Post    `json:"posts,omitempty"`
	Comments     []models.Comment `json:"comments,omitempty"`
}

// FindUserByUsername finds a user by their username
func (s *UserService) FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to retrieve user")
	}

	return &user, nil
}

// GetUserPostCount gets the count of posts authored by a user
func (s *UserService) GetUserPostCount(userID string) int64 {
	var count int64
	database.DB.Model(&models.Post{}).Where("author_id = ?", userID).Count(&count)
	return count
}

// GetUserCommentCount gets the count of comments authored by a user
func (s *UserService) GetUserCommentCount(userID string) int64 {
	var count int64
	database.DB.Model(&models.Comment{}).Where("author_id = ?", userID).Count(&count)
	return count
}

// GetUserPosts retrieves all posts authored by a user
func (s *UserService) GetUserPosts(userID string) ([]models.Post, error) {
	var posts []models.Post
	err := database.DB.
		Select("id", "title", "content", "topic_id", "author_id", "is_pinned", "created_at", "updated_at", "image_url").
		Preload("Topic").
		Where("author_id = ?", userID).
		Order("created_at DESC").
		Find(&posts).Error

	if err != nil {
		return nil, errors.New("failed to retrieve user posts")
	}

	return posts, nil
}

// GetUserComments retrieves all comments authored by a user
func (s *UserService) GetUserComments(userID string) ([]models.Comment, error) {
	var comments []models.Comment
	err := database.DB.
		Select("id", "post_id", "author_id", "content", "is_pinned", "created_at", "updated_at").
		Where("author_id = ?", userID).
		Order("created_at DESC").
		Find(&comments).Error

	if err != nil {
		return nil, errors.New("failed to retrieve user comments")
	}

	return comments, nil
}

// GetUserProfile builds a complete user profile with optional posts and comments
func (s *UserService) GetUserProfile(username string, includePosts bool, includeComments bool) (*UserProfile, error) {
	// Find user by username
	user, err := s.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}

	// Build profile with counts
	profile := &UserProfile{
		User:         *user,
		PostCount:    s.GetUserPostCount(user.ID),
		CommentCount: s.GetUserCommentCount(user.ID),
	}

	// Optionally include posts
	if includePosts {
		posts, err := s.GetUserPosts(user.ID)
		if err != nil {
			return nil, err
		}
		profile.Posts = posts
	}

	// Optionally include comments
	if includeComments {
		comments, err := s.GetUserComments(user.ID)
		if err != nil {
			return nil, err
		}
		profile.Comments = comments
	}

	return profile, nil
}