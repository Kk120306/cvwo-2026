package controllers

import (
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

// func that allows user to like/dislike a post or comment
func CreateOrUpdateVote(c *gin.Context) {
	var body VoteRequest

	// Bind request body to see if all fields are present
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Use Middleware to assume user is authenticated
	user := c.MustGet("user").(models.User)

	// Ensure its eithe a post or comment
	if body.VotableType != "post" && body.VotableType != "comment" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content"})
		return
	}

	// ensure its like or dislike
	if body.VoteType != "like" && body.VoteType != "dislike" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote type"})
		return
	}

	// Try to find existing vote
	var vote models.Vote
	res := database.DB.Where("user_id = ? AND votable_id = ? AND votable_type = ?", user.ID, body.VotableID, body.VotableType).First(&vote)

	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			// No vote exists, create new vote
			newVote := models.Vote{
				UserID:      user.ID,
				VotableID:   body.VotableID,
				VotableType: body.VotableType,
				VoteType:    body.VoteType,
			}
			// Creating the vote into the database
			voteCreated := database.DB.Create(&newVote)
			if voteCreated.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
				return
			}

			// Vote has been created
			c.JSON(http.StatusOK, gin.H{"message": "Vote created", "vote": newVote})
			return
		}

		// Some other DB error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Vote already exists, update vote type if different
	if vote.VoteType != body.VoteType {
		// changing vote in the database
		vote.VoteType = body.VoteType
		changeVote := database.DB.Save(&vote)
		if changeVote.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vote"})
			return
		}
		// Vote has been changed to the new type
		c.JSON(http.StatusOK, gin.H{"message": "Vote updated", "vote": vote})
		return
	}

	// If same vote is clicked, vote is removed
	voteRemove := database.DB.Delete(&vote)
	if voteRemove.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove vote"})
		return
	}

	// Vote has been removed
	c.JSON(http.StatusOK, gin.H{"message": "Vote removed"})
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
