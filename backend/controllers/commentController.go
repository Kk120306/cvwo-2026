package controllers

import (
	"net/http"
	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


// func to get all comments under a certian post
func GetCommentsByPost(c *gin.Context) {
	// Get postID from URL param
	postID := c.Param("postID")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid post ID",
		})
		return
	}

	// Create slice to hold comments
	var comments []models.Comment

	// Query all comments belonging to this post
	// Preload author information so that it can be used to display author details for each comment 
	// https://gorm.io/docs/preload.html
	result := database.DB.
		Where("post_id = ?", postID).
		Preload("Author").
		Order("created_at asc").
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

// func that creates a comment under a certain post 
func CreateComment(c *gin.Context) {
	// Assume user is present from the middlware 
	user := c.MustGet("user").(models.User)

	// Get post ID
	postID := c.Param("postID")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid post ID",
		})
		return
	}

	// Structure of request body
	var body struct {
		Content string 
	}

	// Binding request body failed
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body, please provide valid details",
		})
		return
	}

	// Create comment object
	comment := models.Comment{
		PostID:   postID,
		AuthorID: user.ID,
		Content:  body.Content,
	}

	// Insert into DB
	result := database.DB.Create(&comment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create comment",
		})
		return
	}

	// Return comment
	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
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
		Content string 
	}

	// check body 
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body, please provide valid details",
		})
		return
	}

	// Update the content
	save := database.DB.Model(&comment).Update("content", body.Content)
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

// func to delete a comment, only avaliable for author or admin
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

	// Delete comment
	delete := database.DB.Delete(&comment)
	if delete.Error != nil {
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
