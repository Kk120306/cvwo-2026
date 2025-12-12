package controllers

import (
	"net/http"

	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Function to get posts under a topic
func GetPostsByTopic(c *gin.Context) {
	// getting the slug of topic through params
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid topic slug",
		})
		return
	}

	// Find the topic by slug
	var topic models.Topic
	result := database.DB.First(&topic, "slug = ?", slug)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Topic not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve topic"})
		return
	}

	// Get posts belonging to that topic
	var posts []models.Post
	res := database.DB.
		Where("topic_id = ?", topic.ID).
		Order("is_pinned desc, created_at desc"). // First comes for pinned posts and then we go to order of creation date
		Find(&posts)

		// Error retriveing from db
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve posts",
		})
		return
	}

	// Returning the posts under that topic
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

// Creating a new post under a certain topic - middleware call assumed, that the user is authenticated
func CreatePost(c *gin.Context) {
	// getting the slug of topic through params
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid topic slug",
		})
		return
	}

	// Parse request body
	var body struct {
		Title   string
		Content string
	}

	// check if parsing req binds with struct
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body, please provide valid details",
		})
		return
	}

	// Retrieve topic
	var topic models.Topic
	result := database.DB.First(&topic, "slug = ?", slug)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Topic not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve topic"})
		return
	}

	// Get authenticated user from middleware
	user := c.MustGet("user").(models.User)

	// Create the post
	post := models.Post{
		Title:    body.Title,
		Content:  body.Content,
		TopicID:  topic.ID,
		AuthorID: user.ID,
	}

	// Save to database
	res := database.DB.Create(&post)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create post",
		})
		return
	}

	// Return the created post
	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// function to get a post by their id
func GetPost(c *gin.Context) {
	// getting the id of post through params
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Retrieve post
	var post models.Post
	result := database.DB.First(&post, "id = ?", id)

	// check for any retrival errors
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post"})
		return
	}

	// return the post
	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// funciton to delete a post by their id
func DeletePost(c *gin.Context) {
	// getting the id of post through params
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid post ID"})
		return
	}

	// Fetch the post first
	var post models.Post
	result := database.DB.First(&post, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post"})
		return
	}

	// check if the user is the author or admin
	user := c.MustGet("user").(models.User)
	if post.AuthorID != user.ID && !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this post"})
		return
	}

	// deletion using gorm
	deleteRes := database.DB.Delete(&post)
	if deleteRes.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	// success message
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// Updating a post by id
func UpdatePost(c *gin.Context) {

	// getting the id of post through params
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Parse request body
	var body struct {
		Title   string
		Content string
	}

	// check if parsing req binds with struct
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// Retrieve post
	var post models.Post
	result := database.DB.First(&post, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post"})
		return
	}

	// check if the user is the author or admin
	user := c.MustGet("user").(models.User)
	if post.AuthorID != user.ID && !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this post"})
		return
	}

	// Update fields
	update := database.DB.Model(&post).Updates(models.Post{
		Title:   body.Title,
		Content: body.Content,
	})

	// error during update
	if update.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}
