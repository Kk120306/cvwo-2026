package controllers

import (
	"net/http"

	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/Kk120306/cvwo-2026/backend/services"
	"github.com/gin-gonic/gin"
)

// PostController handles HTTP requests for posts
type PostController struct {
	postService *services.PostService
}

// NewPostController creates a new instance of PostController
func NewPostController() *PostController {
	return &PostController{
		postService: services.NewPostService(),
	}
}

// Function to get all posts (across all topics)
func (pc *PostController) GetAllPosts(c *gin.Context) {
	// Check if user is authenticated
	var userID *string
	u, exists := c.Get("user")
	// if user exists store user data
	if exists {
		user := u.(models.User)
		userID = &user.ID
	}

	// Get posts through service layer
	posts, err := pc.postService.GetAllPosts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// Function to get posts under a topic
func (pc *PostController) GetPostsByTopic(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid topic slug"})
		return
	}

	// Check if user is authenticated
	var userID *string
	if u, exists := c.Get("user"); exists {
		user := u.(models.User)
		userID = &user.ID
	}

	// Get posts through service layer
	posts, err := pc.postService.GetPostsByTopic(slug, userID)
	if err != nil {
		// Determine status code based on error type
		statusCode := http.StatusInternalServerError
		if err.Error() == "topic not found" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// Creating a new post under a certain topic - middleware call assumed, that the user is authenticated
func (pc *PostController) CreatePost(c *gin.Context) {
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
		Title    string  `json:"title" binding:"required"`
		Content  string  `json:"content" binding:"required"`
		ImageUrl *string `json:"imageUrl"`
	}

	// check if parsing req binds with struct
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body, please provide valid details",
		})
		return
	}

	// Retrieve topic through service layer
	topic, err := pc.postService.FindTopicBySlug(slug)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "topic not found" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// Get authenticated user from middleware
	user := c.MustGet("user").(models.User)

	// Create post through service layer
	post, err := pc.postService.CreatePost(services.CreatePostInput{
		Title:    body.Title,
		Content:  body.Content,
		TopicID:  topic.ID,
		AuthorID: user.ID,
		ImageUrl: body.ImageUrl,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the created post
	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// function to get a post by their id
func (pc *PostController) GetPost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Check if user is authenticated
	var userID *string
	if u, exists := c.Get("user"); exists {
		user := u.(models.User)
		userID = &user.ID
	}

	// Get post through service layer
	post, err := pc.postService.GetPostByID(id, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "post not found" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// function to delete a post by their id
func (pc *PostController) DeletePost(c *gin.Context) {
	// getting the id of post through params
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid post ID"})
		return
	}

	// Fetch the post first through service layer
	post, err := pc.postService.FindPostByID(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "post not found" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// check if the user is the author or admin through service layer
	user := c.MustGet("user").(models.User)
	if !pc.postService.CanUserModifyPost(&user, post) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this post"})
		return
	}

	// Delete post through service layer
	err = pc.postService.DeletePost(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// success message
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// UpdatePost updates a post by id
func (pc *PostController) UpdatePost(c *gin.Context) {
	// Get the id of post through params
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Parse request body
	var body struct {
		Title    string  `json:"title" binding:"required"`
		Content  string  `json:"content" binding:"required"`
		ImageURL *string `json:"imageUrl"`
	}

	// Check if parsing req binds with struct
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Retrieve post through service layer
	post, err := pc.postService.FindPostByID(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "post not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// Check if the user is the author or admin
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(models.User)

	// Authorization check through service layer
	if !pc.postService.CanUserModifyPost(&user, post) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this post"})
		return
	}

	// Update post through service layer
	err = pc.postService.UpdatePost(post, services.UpdatePostInput{
		Title:    body.Title,
		Content:  body.Content,
		ImageURL: body.ImageURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload the post with relationships (Author and Topic) through service layer
	updatedPost, err := pc.postService.ReloadPostWithRelationships(post.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": updatedPost,
	})
}

// function that toggles post pins / assumed that user is a admin through middleware
func (pc *PostController) TogglePinPost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Parse request body to get  pin state
	var body struct {
		IsPinned bool `json:"isPinned"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Retrieve post through service layer
	post, err := pc.postService.FindPostByID(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "post not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// Update pin status through service layer
	err = pc.postService.TogglePinPost(post, body.IsPinned)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"isPinned": post.IsPinned,
	})
}
