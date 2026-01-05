package services

import (
	"errors"
	"strings"

	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
)

// PostService handles post business logic
type PostService struct{}

// NewPostService creates a new instance of PostService
func NewPostService() *PostService {
	return &PostService{}
}

// PostWithVotes represents a post with vote counts
type PostWithVotes struct {
	models.Post
	Likes    int64   `json:"likes"`
	Dislikes int64   `json:"dislikes"`
	MyVote   *string `json:"myVote,omitempty"`
}

// CreatePostInput represents the data needed to create a post
type CreatePostInput struct {
	Title    string
	Content  string
	TopicID  string
	AuthorID string
	ImageUrl *string
}

// UpdatePostInput represents the data needed to update a post
type UpdatePostInput struct {
	Title    string
	Content  string
	ImageURL *string
}

// GetAllPosts retrieves all posts across all topics with vote counts
func (s *PostService) GetAllPosts(userID *string) ([]PostWithVotes, error) {
	var posts []PostWithVotes
	var joinUserVote bool

	// if user exists we will join user votes
	if userID != nil && *userID != "" {
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
		`, *userID)
	}

	// Groups by post id, removes any duplicate from join post and combines rows
	// Prioritize pinned post first
	query = query.Preload("Author").
		Preload("Topic").
		Group("posts.id").
		Order("is_pinned DESC, created_at DESC")

	err := query.Find(&posts).Error
	if err != nil {
		return nil, errors.New("failed to retrieve posts")
	}

	return posts, nil
}

// GetPostsByTopic retrieves all posts under a specific topic with vote counts
func (s *PostService) GetPostsByTopic(slug string, userID *string) ([]PostWithVotes, error) {
	// Normalize slug
	slug = strings.ToLower(slug)

	// Find topic first
	var topic models.Topic
	err := database.DB.First(&topic, "slug = ?", slug).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("topic not found")
		}
		return nil, errors.New("failed to retrieve topic")
	}

	var posts []PostWithVotes
	var joinUserVote bool

	// Check if user is authenticated, if so store user data
	if userID != nil && *userID != "" {
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
		`, *userID)
	}

	// Groups by post id, removes any duplicate from join post and combines rows
	query = query.Preload("Author").
		Preload("Topic").
		Group("posts.id").
		Order("is_pinned DESC, created_at DESC")

	// Execute query and put results in posts slice
	findErr := query.Find(&posts).Error
	if findErr != nil {
		return nil, errors.New("failed to retrieve posts")
	}

	return posts, nil
}

// CreatePost creates a new post under a topic
func (s *PostService) CreatePost(input CreatePostInput) (*models.Post, error) {
	// Validate input
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Content) == "" {
		return nil, errors.New("title and content cannot be empty")
	}

	// Sanitize content input from rich text editor
	// https://github.com/microcosm-cc/bluemonday - prevent xss attacks
	safeContent := bluemonday.UGCPolicy().Sanitize(input.Content)

	// Create the post
	post := models.Post{
		Title:    input.Title,
		Content:  safeContent,
		TopicID:  input.TopicID,
		AuthorID: input.AuthorID,
		ImageUrl: input.ImageUrl,
	}

	// Save to database
	res := database.DB.Create(&post)
	if res.Error != nil {
		return nil, errors.New("failed to create post")
	}

	return &post, nil
}

// GetPostByID retrieves a single post by ID with vote counts
func (s *PostService) GetPostByID(id string, userID *string) (*PostWithVotes, error) {
	var post PostWithVotes
	var joinUserVote bool

	// if user exists store user data
	if userID != nil && *userID != "" {
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
		`, *userID)
	}

	// Groups by post id, removes any duplicate from join post and combines rows
	// Preload Author and Topic relationships
	query = query.Preload("Author").
		Preload("Topic").
		Group("posts.id")

	err := query.First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, errors.New("failed to retrieve post")
	}

	return &post, nil
}

// FindPostByID finds a post by ID (without vote counts)
func (s *PostService) FindPostByID(id string) (*models.Post, error) {
	var post models.Post
	result := database.DB.First(&post, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("post not found")
		}
		return nil, errors.New("failed to retrieve post")
	}
	return &post, nil
}

// UpdatePost updates a post's content
func (s *PostService) UpdatePost(post *models.Post, input UpdatePostInput) error {
	// Validate fields are not empty
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Content) == "" {
		return errors.New("title and content cannot be empty")
	}

	// Sanitize content
	safeContent := bluemonday.UGCPolicy().Sanitize(input.Content)

	// Update fields - Use map to allow nil values
	updates := map[string]interface{}{
		"title":   input.Title,
		"content": safeContent,
	}

	// Handle ImageURL separately to allow setting it to null
	if input.ImageURL != nil {
		updates["image_url"] = *input.ImageURL
	} else {
		// Explicitly set to null if imageUrl is null
		updates["image_url"] = nil
	}

	updateErr := database.DB.Model(post).Updates(updates).Error
	if updateErr != nil {
		return errors.New("failed to update post")
	}

	return nil
}

// DeletePost deletes a post and all associated data (votes, comments, comment votes)
func (s *PostService) DeletePost(post *models.Post) error {
	// Transaction to handle all votes, ensures all or no operations happen
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Get all comment IDs for this post to delete their votes
		var commentIDs []string
		retrieveErr := tx.Model(&models.Comment{}).
			Where("post_id = ?", post.ID). // Get comments under the post
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
		err := tx.Where("votable_id = ? AND votable_type = ?", post.ID, "post").
			Delete(&models.Vote{}).Error
		if err != nil {
			return err
		}

		// Delete comments on the post
		commentErr := tx.Where("post_id = ?", post.ID).Delete(&models.Comment{}).Error
		if commentErr != nil {
			return commentErr
		}

		// Finally, delete the post itself
		delError := tx.Delete(post).Error
		if delError != nil {
			return delError
		}

		return nil
	})

	if err != nil {
		return errors.New("failed to delete post")
	}

	return nil
}

// TogglePinPost toggles the pin status of a post
func (s *PostService) TogglePinPost(post *models.Post, isPinned bool) error {
	// Update pin status
	post.IsPinned = isPinned
	saveErr := database.DB.Save(post).Error
	if saveErr != nil {
		return errors.New("failed to update pin status")
	}

	return nil
}

// CanUserModifyPost checks if a user has permission to modify a post
func (s *PostService) CanUserModifyPost(user *models.User, post *models.Post) bool {
	// check if user is permitted (either author or admin)
	return user.ID == post.AuthorID || user.IsAdmin
}

// FindTopicBySlug finds a topic by its slug
func (s *PostService) FindTopicBySlug(slug string) (*models.Topic, error) {
	slug = strings.ToLower(slug)

	var topic models.Topic
	result := database.DB.First(&topic, "slug = ?", slug)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("topic not found")
		}
		return nil, errors.New("failed to retrieve topic")
	}
	return &topic, nil
}

// ReloadPostWithRelationships reloads a post with Author and Topic relationships
func (s *PostService) ReloadPostWithRelationships(postID string) (*models.Post, error) {
	var post models.Post
	reloadErr := database.DB.
		Preload("Author").
		Preload("Topic").
		Where("id = ?", postID).
		First(&post).Error
	if reloadErr != nil {
		return nil, errors.New("failed to fetch updated post")
	}
	return &post, nil
}
