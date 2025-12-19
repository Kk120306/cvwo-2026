package controllers

import (
	"net/http"
	"strings"

	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
)

// Struct for comments with vote counts
type CommentWithVotes struct {
	models.Comment
	Likes    int64
	Dislikes int64
	MyVote   *string
}

// func to get all comments under a certian post
func GetCommentsByPost(c *gin.Context) {
	// Get postID from URL param
	postID := c.Param("postId")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid post ID",
		})
		return
	}

	// Create slice to hold comments
	var comments []CommentWithVotes

	// Check if user is authenticated
	var userID string
	var joinUserVote bool
	u, exists := c.Get("user")
	// If user is authenticated - we remember userID to join later
	if exists {
		user := u.(models.User)
		userID = user.ID
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
	// THen use selectStr to get the counts. Since LEFT JOIN is used, Coalesce is used in select str
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
		`, userID)
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve comments",
		})
		return
	}

	// Return the comments
	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}

// CreateComment creates a comment under a certain post
func CreateComment(c *gin.Context) {
	// Get user from middleware with error handling
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(models.User)

	// Get and validate post ID
	postID := c.Param("postId")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid post ID"})
		return
	}

	// Check if post exists
	var post models.Post
	if err := database.DB.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var body struct {
		Content  string
		ImageUrl *string
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Validate content is not empty after trimming
	if strings.TrimSpace(body.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content cannot be empty"})
		return
	}

	// Sanitize content from the rich text editor
	// https://github.com/microcosm-cc/bluemonday - prevent XSS attacks
	safeContent := bluemonday.UGCPolicy().Sanitize(body.Content)

	// Create comment object
	comment := models.Comment{
		PostID:   postID,
		AuthorID: user.ID,
		Content:  safeContent,
		ImageUrl: body.ImageUrl,
	}

	// Insert into DB
	createErr := database.DB.Create(&comment).Error
	if createErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Fetch the created comment with author
	fetchErr := database.DB.Preload("Author").Where("id = ?", comment.ID).First(&comment).Error
	if fetchErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment"})
		return
	}

	// We do this cus retriving votes is more work and also it has been normlaized to the correct json structure of a Comment
	// Since thats the only use case of ceating a comment
	response := gin.H{
		"comment": gin.H{
			"id":        comment.ID,
			"postId":    comment.PostID,
			"authorId":  comment.AuthorID,
			"imageUrl":  comment.ImageUrl,
			"content":   comment.Content,
			"createdAt": comment.CreatedAt,
			"updatedAt": comment.UpdatedAt,
			"author": gin.H{
				"id":        comment.Author.ID,
				"username":  comment.Author.Username,
				"avatarURL": comment.Author.AvatarURL,
				"isAdmin":   comment.Author.IsAdmin,
			},
			"likes":    0,
			"dislikes": 0,
			"myVote":   nil,
			"isPinned": comment.IsPinned,
		},
	}

	c.JSON(http.StatusOK, response)
}

// Function to update comments only accessible by either the author or an admin
func UpdateComment(c *gin.Context) {
	// Get authenticated user from middleware
	user := c.MustGet("user").(models.User)

	// Get comment ID
	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid comment ID",
		})
		return
	}

	// Find the comment in DB
	var comment models.Comment
	result := database.DB.First(&comment, "id = ?", commentID)

	// If not found
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Comment not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve comment",
		})
		return
	}

	// check if user is permitted
	if user.ID != comment.AuthorID && !user.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "You are not allowed to update this comment",
		})
		return
	}

	// Parse body
	var body struct {
		Content  string
		ImageURL *string
	}

	// check body
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body, please provide valid details",
		})
		return
	}
	// updates map
	updates := map[string]interface{}{
		"content": body.Content,
	}

	// Add ImageURL to updates if provided
	if body.ImageURL != nil {
		updates["image_url"] = *body.ImageURL
	}

	// Update the comment with all fields
	save := database.DB.Model(&comment).Updates(updates)
	if save.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update comment",
		})
		return
	}

	// Return updated comment
	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
}

// func to delete a comment, only available for author or admin
func DeleteComment(c *gin.Context) {
	// Get user from middleware
	user := c.MustGet("user").(models.User)

	// Get comment ID
	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid comment ID",
		})
		return
	}

	// Find comment
	var comment models.Comment
	result := database.DB.First(&comment, "id = ?", commentID)

	// If not found
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Comment not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve comment",
		})
		return
	}

	// check if permitted
	if user.ID != comment.AuthorID && !user.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "You are not allowed to delete this comment",
		})
		return
	}

	// Transaction to delete votes and comment - ensures that every opeartion happns or none at all
	// https://gorm.io/docs/transactions.html
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Delete all votes on this comment (polymorphic relationship)
		err := tx.Where("votable_id = ? AND votable_type = ?", commentID, "comment").Delete(&models.Vote{}).Error
		if err != nil {
			return err
		}

		// 2. Delete the comment itself
		delErr := tx.Delete(&comment).Error
		if delErr != nil {
			return delErr
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete comment",
		})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}
