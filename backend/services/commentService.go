package services

import (
	"errors"
	"strings"

	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
)

// CommentService handles comment business logic
type CommentService struct{}

// NewCommentService creates a new instance of CommentService
func NewCommentService() *CommentService {
	return &CommentService{}
}

// CommentWithVotes represents a comment with vote counts
type CommentWithVotes struct {
	models.Comment
	Likes    int64   `json:"likes"`
	Dislikes int64   `json:"dislikes"`
	MyVote   *string `json:"myVote,omitempty"`
}

// CreateCommentInput represents the data needed to create a comment
type CreateCommentInput struct {
	PostID   string
	AuthorID string
	Content  string
}

// UpdateCommentInput represents the data needed to update a comment
type UpdateCommentInput struct {
	Content string
}

// TogglePinInput represents the data needed to toggle pin status
type TogglePinInput struct {
	IsPinned bool
	AuthorID string
}

// GetCommentsByPost retrieves all comments for a specific post with vote counts
func (s *CommentService) GetCommentsByPost(postID string, userID *string) ([]CommentWithVotes, error) {
	// Create slice to hold comments
	var comments []CommentWithVotes

	// Check if user is authenticated to join user votes
	var joinUserVote bool
	if userID != nil && *userID != "" {
		joinUserVote = true
	}

	// Even if user is not logged in we still want the likes and dislikes count for each comment
	selectStr := `
		comments.*,
		COALESCE(SUM(CASE WHEN votes.vote_type = 'like' THEN 1 ELSE 0 END), 0) AS likes,
		COALESCE(SUM(CASE WHEN votes.vote_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes
	`

	// If user is logged in, we add the my_vote field which is used to show what vote the user has made
	if joinUserVote {
		selectStr += `,
			MAX(user_votes.vote_type) AS my_vote`
	}

	// The query to pass
	// For comments that belong to the postID, attach all votes that belong to each comment
	// Then use selectStr to get the counts. Since LEFT JOIN is used, Coalesce is used in select str
	query := database.DB.
		Model(&models.Comment{}).
		Select(selectStr).
		Joins(`
			LEFT JOIN votes 
			ON votes.votable_id = comments.id 
			AND votes.votable_type = 'comment'
		`).
		Where("comments.post_id = ?", postID)

	// if user is logged in, combine a second join where we only want to join votes that belong to the user
	if joinUserVote {
		query = query.Joins(`
			LEFT JOIN votes AS user_votes
			ON user_votes.votable_id = comments.id
			AND user_votes.votable_type = 'comment'
			AND user_votes.user_id = ?
		`, *userID)
	}

	// Execute query
	// Get any details about Author of comment
	// Group so that there are no duplicate comments from Join
	// So that for each comment we get the total and not each vote row
	// Also for repeated in user_votes if user exists
	result := query.
		Preload("Author").
		Group("comments.id").
		Order("comments.created_at asc").
		Find(&comments)

	// If database error
	if result.Error != nil {
		return nil, errors.New("failed to retrieve comments")
	}

	return comments, nil
}

// CreateComment creates a new comment under a post
func (s *CommentService) CreateComment(input CreateCommentInput) (*models.Comment, error) {
	// Validate content is not empty after trimming
	if strings.TrimSpace(input.Content) == "" {
		return nil, errors.New("content cannot be empty")
	}

	// Sanitize content from the rich text editor
	// https://github.com/microcosm-cc/bluemonday - prevent XSS attacks
	safeContent := bluemonday.UGCPolicy().Sanitize(input.Content)

	// Create comment object
	comment := models.Comment{
		PostID:   input.PostID,
		AuthorID: input.AuthorID,
		Content:  safeContent,
	}

	// Insert into DB
	createErr := database.DB.Create(&comment).Error
	if createErr != nil {
		return nil, errors.New("failed to create comment")
	}

	// Fetch the created comment with author
	fetchErr := database.DB.Preload("Author").Where("id = ?", comment.ID).First(&comment).Error
	if fetchErr != nil {
		return nil, errors.New("failed to fetch comment")
	}

	return &comment, nil
}

// FindCommentByID finds a comment by its ID
func (s *CommentService) FindCommentByID(commentID string) (*models.Comment, error) {
	// Find the comment in DB
	var comment models.Comment
	result := database.DB.First(&comment, "id = ?", commentID)

	// If not found
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("comment not found")
		}
		return nil, errors.New("failed to retrieve comment")
	}

	return &comment, nil
}

// UpdateComment updates a comment's content
func (s *CommentService) UpdateComment(comment *models.Comment, input UpdateCommentInput) error {
	// Validate content is not empty after trimming
	if strings.TrimSpace(input.Content) == "" {
		return errors.New("content cannot be empty")
	}

	// Sanitize content
	safeContent := bluemonday.UGCPolicy().Sanitize(input.Content)

	// Update the content
	save := database.DB.Model(comment).Update("content", safeContent)
	if save.Error != nil {
		return errors.New("failed to update comment")
	}

	return nil
}

// DeleteComment deletes a comment and all its votes
func (s *CommentService) DeleteComment(comment *models.Comment) error {
	// Transaction to delete votes and comment - ensures that every operation happens or none at all
	// https://gorm.io/docs/transactions.html
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Delete all votes on this comment (polymorphic relationship)
		err := tx.Where("votable_id = ? AND votable_type = ?", comment.ID, "comment").Delete(&models.Vote{}).Error
		if err != nil {
			return err
		}

		// 2. Delete the comment itself
		delErr := tx.Delete(comment).Error
		if delErr != nil {
			return delErr
		}

		return nil
	})

	if err != nil {
		return errors.New("failed to delete comment")
	}

	return nil
}

// CanUserModifyComment checks if a user has permission to modify a comment
func (s *CommentService) CanUserModifyComment(user *models.User, comment *models.Comment) bool {
	// check if user is permitted (either author or admin)
	return user.ID == comment.AuthorID || user.IsAdmin
}

// PostExists checks if a post exists by ID
func (s *CommentService) PostExists(postID string) (bool, error) {
	// Check if post exists
	var post models.Post
	if err := database.DB.First(&post, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, errors.New("database error")
	}
	return true, nil
}
