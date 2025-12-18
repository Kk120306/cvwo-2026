package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
)

// Struct to hold vote data along with post
type PostWithVotes struct {
	models.Post
	Likes    int64
	Dislikes int64
	MyVote   *string
}

// Function to get all posts (across all topics)
func GetAllPosts(c *gin.Context) {
	var posts []PostWithVotes
	var userID string
	var joinUserVote bool

	// if user exists store user data
	u, exists := c.Get("user")
	if exists {
		user := u.(models.User)
		userID = user.ID
		joinUserVote = true
	}

	// From all posts get the likes and dislikes count
	selectStr := `
		posts.*,
		COALESCE(SUM(CASE WHEN votes.vote_type = 'like' THEN 1 ELSE 0 END), 0) AS likes,
		COALESCE(SUM(CASE WHEN votes.vote_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes
	`

	// If user is authenticated, also get their vote on each post
	if joinUserVote {
		selectStr += `,
			MAX(user_votes.vote_type) AS my_vote`
	}

	// attach votes table where votable_id matches post id and votable_type is post
	// left joins creates null so coalesce is needed
	query := database.DB.Model(&models.Post{}).
		Select(selectStr).
		Joins(`
			LEFT JOIN votes 
			ON votes.votable_id = posts.id 
			AND votes.votable_type = 'post'
		`)

	// if user is authenticated, join another table that only has users votes
	if joinUserVote {
		query = query.Joins(`
			LEFT JOIN votes AS user_votes
			ON user_votes.votable_id = posts.id
			AND user_votes.votable_type = 'post'
			AND user_votes.user_id = ?
		`, userID)
	}

	// Groups by post id, reomves any duplicate from join post and combines rows
	// Prioritize pinned post first
	query = query.Preload("Author").
		Preload("Topic").
		Group("posts.id").
		Order("is_pinned DESC, created_at DESC")

	err := query.Find(&posts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// Function to get posts under a topic
func GetPostsByTopic(c *gin.Context) {
	slug := strings.ToLower(c.Param("slug"))
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid topic slug"})
		return
	}

	var topic models.Topic
	err := database.DB.First(&topic, "slug = ?", slug).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Topic not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve topic"})
		return
	}

	var posts []PostWithVotes
	var userID string
	var joinUserVote bool
	// Check if user is authenticated, if so store user data
	if u, exists := c.Get("user"); exists {
		user := u.(models.User)
		userID = user.ID
		joinUserVote = true
	}

	// From all posts get the likes and dislikes count
	selectStr := `
		posts.*,
		COALESCE(SUM(CASE WHEN votes.vote_type = 'like' THEN 1 ELSE 0 END), 0) AS likes,
		COALESCE(SUM(CASE WHEN votes.vote_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes
	`

	// If user is authenticated, also get their vote on each post
	if joinUserVote {
		selectStr += `,
			MAX(CASE WHEN user_votes.user_id IS NOT NULL THEN user_votes.vote_type ELSE NULL END) AS my_vote`
	}

	// attach votes table where votable_id matches post id and votable_type is post
	// left joins creates null so coalesce is needed
	query := database.DB.Model(&models.Post{}).
		Select(selectStr).
		Joins(`
			LEFT JOIN votes 
			ON votes.votable_id = posts.id 
			AND votes.votable_type = 'post'
		`).
		Where("posts.topic_id = ?", topic.ID)

	// if user is authenticated, join another table that only has users votes
	if joinUserVote {
		query = query.Joins(`
			LEFT JOIN votes AS user_votes
			ON user_votes.votable_id = posts.id
			AND user_votes.votable_type = 'post'
			AND user_votes.user_id = ?
		`, userID)
	}

	// Groups by post id, reomves any duplicate from join post and combines rows
	query = query.Preload("Author").
		Preload("Topic").
		Group("posts.id").
		Order("is_pinned DESC, created_at DESC")

	// Execute query and put results in posts slice
	findErr := query.Find(&posts).Error
	if findErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
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
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Topic not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve topic"})
		return
	}

	// Get authenticated user from middleware
	user := c.MustGet("user").(models.User)

	// Sanitize content input from rich text editor
	// https://github.com/microcosm-cc/bluemonday - prevent xxs attacks
	safeContent := bluemonday.UGCPolicy().Sanitize(body.Content)

	// Create the post
	post := models.Post{
		Title:    body.Title,
		Content:  safeContent,
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
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post PostWithVotes
	var userID string
	var joinUserVote bool

	// if user exists store user data
	if u, exists := c.Get("user"); exists {
		user := u.(models.User)
		userID = user.ID
		joinUserVote = true
	}

	// From all posts get the likes and dislikes count
	selectStr := `
		posts.*,
		COALESCE(SUM(CASE WHEN votes.vote_type = 'like' THEN 1 ELSE 0 END), 0) AS likes,
		COALESCE(SUM(CASE WHEN votes.vote_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes
	`

	// If user is authenticated, also get their vote on the post
	if joinUserVote {
		selectStr += `,
			MAX(CASE WHEN user_votes.user_id IS NOT NULL THEN user_votes.vote_type ELSE NULL END) AS my_vote`
	}

	// attach votes table where votable_id matches post id and votable_type is post
	// left joins creates null so coalesce is needed
	query := database.DB.Model(&models.Post{}).
		Select(selectStr).
		Joins(`
			LEFT JOIN votes 
			ON votes.votable_id = posts.id 
			AND votes.votable_type = 'post'
		`).
		Where("posts.id = ?", id)

	// if user is authenticated, join another table that only has users votes
	if joinUserVote {
		query = query.Joins(`
			LEFT JOIN votes AS user_votes
			ON user_votes.votable_id = posts.id
			AND user_votes.votable_type = 'post'
			AND user_votes.user_id = ?
		`, userID)
	}

	// Groups by post id, reomves any duplicate from join post and combines rows
	// Preload Author and Topic relationships
	query = query.Preload("Author").
		Preload("Topic").
		Group("posts.id")

	err := query.First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// function to delete a post by their id
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

	// Transaction to handle all votes, ensures all or no operations happen
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Get all comment IDs for this post to delete their votes
		var commentIDs []string
		retrieveErr := tx.Model(&models.Comment{}).
			Where("post_id = ?", id).      //  Get comments under the post
			Pluck("id", &commentIDs).Error // Pluck gets a slice of only the ids
		if retrieveErr != nil {
			return retrieveErr
		}

		// Delete votes on comments
		if len(commentIDs) > 0 {
			err := tx.Where("votable_id IN ? AND votable_type = ?", commentIDs, "comment").
				Delete(&models.Vote{}).Error
			if err != nil {
				return err
			}
		}

		// Delete votes on the post itself
		err := tx.Where("votable_id = ? AND votable_type = ?", id, "post").
			Delete(&models.Vote{}).Error
		if err != nil {
			return err
		}

		//  Delete comments on the post
		commentErr := tx.Where("post_id = ?", id).Delete(&models.Comment{}).Error
		if commentErr != nil {
			return commentErr
		}

		// Finally, delete the post itself
		delError := tx.Delete(&post).Error
		if delError != nil {
			return delError
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	// success message
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// UpdatePost updates a post by id
func UpdatePost(c *gin.Context) {
	// Get the id of post through params
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Parse request body
	var body struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	// Check if parsing req binds with struct
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Validate fields are not empty
	if strings.TrimSpace(body.Title) == "" || strings.TrimSpace(body.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and content cannot be empty"})
		return
	}

	// Retrieve post
	var post models.Post
	retriveErr := database.DB.First(&post, "id = ?", id).Error
	if retriveErr != nil {
		if errors.Is(retriveErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post"})
		return
	}

	// Check if the user is the author or admin
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(models.User)

	// Authorization check
	if post.AuthorID != user.ID && !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this post"})
		return
	}

	// Update fields
	updateErr := database.DB.Model(&post).Updates(models.Post{
		Title:   body.Title,
		Content: body.Content,
	}).Error
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	// Reload the post with relationships (Author and Topic)
	var updatedPost models.Post
	reloadErr := database.DB.
		Preload("Author").
		Preload("Topic").
		Where("id = ?", post.ID).
		First(&updatedPost).Error
	if reloadErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": updatedPost,
	})
}
