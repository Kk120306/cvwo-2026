package controllers

import (
	"errors"
	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// Structure for creating/updating a vote
type VoteRequest struct {
	VotableID   string // The ID of the post or comment being voted on
	VotableType string // "post" or "comment"
	VoteType    string // "like" or "dislike"
}

// CreateOrUpdateVote allows a user to like/dislike a post or comment
func CreateOrUpdateVote(c *gin.Context) {
	var body VoteRequest

	// Bind request body
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get authenticated user
	user := c.MustGet("user").(models.User)

	// Validate votable type
	if body.VotableType != "post" && body.VotableType != "comment" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}

	// Validate vote type
	if body.VoteType != "like" && body.VoteType != "dislike" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote type"})
		return
	}

	// Try to find existing vote
	var vote models.Vote
	err := database.DB.Where("user_id = ? AND votable_id = ? AND votable_type = ?", user.ID, body.VotableID, body.VotableType).First(&vote).Error

	// checking for any errors
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // Use error package cus error can be wrapped
			// No vote exists, create new vote
			newVote := models.Vote{
				UserID:      user.ID,
				VotableID:   body.VotableID,
				VotableType: body.VotableType,
				VoteType:    body.VoteType,
			}
			// Checking if there is a error creating the new vote
			createErr := database.DB.Create(&newVote).Error
			if createErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
				return
			}
		} else {
			// Catching a error on the Database side
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	} else {
		// Vote exists
		if vote.VoteType == body.VoteType {
			// Same vote clicked, remove it
			delErr := database.DB.Delete(&vote).Error
			if delErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove vote"})
				return
			}
		} else {
			// Different vote, update it
			vote.VoteType = body.VoteType
			saveErr := database.DB.Save(&vote).Error
			if saveErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vote"})
				return
			}
		}
	}

	// Respond with updated counts - helper that passes vote info with attachment to user
	respondWithVoteCountsAndUserVote(c, body.VotableID, body.VotableType, user.ID)
}

// Helper: respond with vote counts and user vote using a single query
func respondWithVoteCountsAndUserVote(c *gin.Context, votableID, votableType, userID string) {
	// where the result is stored
	var result struct {
		Likes    int64   `gorm:"column:likes"`
		Dislikes int64   `gorm:"column:dislikes"`
		MyVote   *string `gorm:"column:my_vote"` // use pointer so can tell null from "like"/"dislike"
	}

	// Database query
	database.DB.Model(&models.Vote{}).
		// sum case adds up likes and dislikes when vote_type is met and stored as likes and dislikes
		// Coalesce helps to convert any null rows to 0
		// For each row keep only if user id matches, else null and then max just removes all nulls
		Select(`
			COALESCE(SUM(CASE WHEN vote_type = 'like' THEN 1 ELSE 0 END),0) AS likes,
			COALESCE(SUM(CASE WHEN vote_type = 'dislike' THEN 1 ELSE 0 END),0) AS dislikes,
			MAX(CASE WHEN user_id = ? THEN vote_type ELSE NULL END) AS my_vote
		`, userID).
		Where("votable_id = ? AND votable_type = ?", votableID, votableType).
		Scan(&result)

	c.JSON(http.StatusOK, gin.H{
		"likes":    result.Likes,
		"dislikes": result.Dislikes,
		"myVote":   result.MyVote,
	})
}

// func that returns the total votes of a certain post or comment
// Using params as its a func that retrives info for the client
func GetVotesCount(c *gin.Context) {
	votableID := c.Param("id")
	votableType := c.Query("type") // must be "post" or "comment"

	// Check if params have valid values
	if votableID == "" || (votableType != "post" && votableType != "comment") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid votable ID or type"})
		return
	}

	// varaibles to hold the amount of likes and dislikes
	var likes int64
	var dislikes int64

	// Count likes from the database
	likeQuery := database.DB.Model(&models.Vote{}).Where("votable_id = ? AND votable_type = ? AND vote_type = ?", votableID, votableType, "like").Count(&likes)
	if likeQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Count dislikes in the database
	dislikeQuery := database.DB.Model(&models.Vote{}).Where("votable_id = ? AND votable_type = ? AND vote_type = ?", votableID, votableType, "dislike").Count(&dislikes)
	if dislikeQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Return the counts as JSON response
	c.JSON(http.StatusOK, gin.H{
		"votable_id":   votableID,
		"votable_type": votableType,
		"likes":        likes,
		"dislikes":     dislikes,
	})
}
