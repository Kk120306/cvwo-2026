package controllers

import (
	"net/http"

	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/Kk120306/cvwo-2026/backend/services"
	"github.com/gin-gonic/gin"
)

// VoteController handles HTTP requests for votes
type VoteController struct {
	voteService *services.VoteService
}

// NewVoteController creates a new instance of VoteController
func NewVoteController() *VoteController {
	return &VoteController{
		voteService: services.NewVoteService(),
	}
}

// Structure for creating/updating a vote
type VoteRequest struct {
	VotableID   string `json:"votableId" binding:"required"`   // ID of the post or comment
	VotableType string `json:"votableType" binding:"required"` // "post" or "comment"
	VoteType    string `json:"voteType" binding:"required"`    // "like" or "dislike"
}

// CreateOrUpdateVote allows a user to like/dislike a post or comment
func (vc *VoteController) CreateOrUpdateVote(c *gin.Context) {
	var body VoteRequest
	// Bind request body
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get authenticated user
	user := c.MustGet("user").(models.User)

	// Validate input through service layer
	input := services.VoteInput{
		VotableID:   body.VotableID,
		VotableType: body.VotableType,
		VoteType:    body.VoteType,
	}

	err := vc.voteService.ValidateVoteInput(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Try to find existing vote through service layer
	vote, err := vc.voteService.FindExistingVote(user.ID, body.VotableID, body.VotableType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if vote == nil {
		// No vote exists, create new vote through service layer
		err = vc.voteService.CreateVote(user.ID, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Vote exists
		if vote.VoteType == body.VoteType {
			// Same vote clicked, remove it through service layer
			err = vc.voteService.DeleteVote(vote)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			// Different vote, update it through service layer
			err = vc.voteService.UpdateVote(vote, body.VoteType)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	// Respond with updated counts through service layer
	voteCounts, err := vc.voteService.GetVoteCountsWithUserVote(body.VotableID, body.VotableType, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"likes":    voteCounts.Likes,
		"dislikes": voteCounts.Dislikes,
		"myVote":   voteCounts.MyVote,
	})
}

// func that returns the total votes of a certain post or comment
// Using params as its a func that retrieves info for the client
func (vc *VoteController) GetVotesCount(c *gin.Context) {
	votableID := c.Param("id")
	votableType := c.Query("type") // must be "post" or "comment"

	// Check if params have valid values
	if votableID == "" || (votableType != "post" && votableType != "comment") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid votable ID or type"})
		return
	}

	// Get vote counts through service layer
	likes, dislikes, err := vc.voteService.GetVoteCounts(votableID, votableType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the counts as JSON response
	c.JSON(http.StatusOK, gin.H{
		"votableId":   votableID,
		"votableType": votableType,
		"likes":       likes,
		"dislikes":    dislikes,
	})
}
