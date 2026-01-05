package controllers

import (
	"net/http"

	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/Kk120306/cvwo-2026/backend/services"
	"github.com/gin-gonic/gin"
)

// CommentController handles HTTP requests for comments
type CommentController struct {
	commentService *services.CommentService
}

// NewCommentController creates a new instance of CommentController
func NewCommentController() *CommentController {
	return &CommentController{
		commentService: services.NewCommentService(),
	}
}

// func to get all comments under a certain post
func (cc *CommentController) GetCommentsByPost(c *gin.Context) {
	// Get postID from URL param
	postID := c.Param("postId")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid post ID",
		})
		return
	}

	// Check if user is authenticated
	var userID *string
	u, exists := c.Get("user")
	// If user is authenticated - we remember userID to join later
	if exists {
		user := u.(models.User)
		userID = &user.ID
	}

	// Get comments through service layer
	comments, err := cc.commentService.GetCommentsByPost(postID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the comments
	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}

// CreateComment creates a comment under a certain post
func (cc *CommentController) CreateComment(c *gin.Context) {
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

	// Check if post exists through service layer
	postExists, err := cc.commentService.PostExists(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if !postExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var body struct {
		Content string `json:"content" binding:"required"`
	}

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Create comment through service layer
	comment, err := cc.commentService.CreateComment(services.CreateCommentInput{
		PostID:   postID,
		AuthorID: user.ID,
		Content:  body.Content,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// We do this cus retrieving votes is more work and also it has been normalized to the correct json structure of a Comment
	// Since thats the only use case of creating a comment
	response := gin.H{
		"comment": gin.H{
			"id":        comment.ID,
			"postId":    comment.PostID,
			"authorId":  comment.AuthorID,
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
func (cc *CommentController) UpdateComment(c *gin.Context) {
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

	// Find comment through service layer
	comment, err := cc.commentService.FindCommentByID(commentID)
	if err != nil {
		// Determine status code based on error type
		statusCode := http.StatusInternalServerError
		if err.Error() == "comment not found" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check if user is permitted through service layer
	if !cc.commentService.CanUserModifyComment(&user, comment) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "You are not allowed to update this comment",
		})
		return
	}

	// Parse body
	var body struct {
		Content string `json:"content" binding:"required"`
	}

	// check body
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body, please provide valid details",
		})
		return
	}

	// Update comment through service layer
	err = cc.commentService.UpdateComment(comment, services.UpdateCommentInput{
		Content: body.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return updated comment
	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
}

// func to delete a comment, only available for author or admin
func (cc *CommentController) DeleteComment(c *gin.Context) {
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

	// Find comment through service layer
	comment, err := cc.commentService.FindCommentByID(commentID)
	if err != nil {
		// Determine status code based on error type
		statusCode := http.StatusInternalServerError
		if err.Error() == "comment not found" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check if permitted through service layer
	if !cc.commentService.CanUserModifyComment(&user, comment) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "You are not allowed to delete this comment",
		})
		return
	}

	// Delete comment through service layer
	err = cc.commentService.DeleteComment(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}

// function that toggles comment pin / assumes the user is logged in
func (cc *CommentController) TogglePinComment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid comment ID",
		})
		return
	}

	var body struct {
		IsPinned bool   `json:"isPinned" binding:"required"`
		AuthorID string `json:"authorId" binding:"required"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// Get the authenticated user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	authenticatedUser := user.(models.User)

	// Find the comment through service layer
	comment, err := cc.commentService.FindCommentByID(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "comment not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if the authenticated user is the author of the comment OR an admin through service layer
	if !cc.commentService.CanUserModifyComment(&authenticatedUser, comment) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only pin/unpin your own comments unless you are an admin",
		})
		return
	}

	// Update the pin status through service layer
	err = cc.commentService.TogglePinComment(comment, body.IsPinned)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
}
