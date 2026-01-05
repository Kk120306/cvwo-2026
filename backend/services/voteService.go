package services

import (
	"errors"

	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"gorm.io/gorm"
)

// VoteService handles vote business logic
type VoteService struct{}

// NewVoteService creates a new instance of VoteService
func NewVoteService() *VoteService {
	return &VoteService{}
}

// VoteInput represents the data needed to create/update a vote
type VoteInput struct {
	VotableID   string
	VotableType string
	VoteType    string
}

// VoteCounts represents vote statistics for a votable item
type VoteCounts struct {
	Likes    int64   `json:"likes"`
	Dislikes int64   `json:"dislikes"`
	MyVote   *string `json:"myVote"`
}

// FindExistingVote finds a user's existing vote on a votable item
func (s *VoteService) FindExistingVote(userID, votableID, votableType string) (*models.Vote, error) {
	var vote models.Vote
	err := database.DB.Where("user_id = ? AND votable_id = ? AND votable_type = ?",
		userID, votableID, votableType).First(&vote).Error

	// checking for any errors
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // Use error package cus error can be wrapped
			return nil, nil // No vote exists (not an error)
		}
		return nil, errors.New("database error")
	}

	return &vote, nil
}

// CreateVote creates a new vote
func (s *VoteService) CreateVote(userID string, input VoteInput) error {
	newVote := models.Vote{
		UserID:      userID,
		VotableID:   input.VotableID,
		VotableType: input.VotableType,
		VoteType:    input.VoteType,
	}

	// Checking if there is a error creating the new vote
	createErr := database.DB.Create(&newVote).Error
	if createErr != nil {
		return errors.New("failed to create vote")
	}

	return nil
}

// DeleteVote deletes an existing vote (when user clicks same vote again)
func (s *VoteService) DeleteVote(vote *models.Vote) error {
	delErr := database.DB.Delete(vote).Error
	if delErr != nil {
		return errors.New("failed to remove vote")
	}
	return nil
}

// UpdateVote updates an existing vote (when user clicks different vote)
func (s *VoteService) UpdateVote(vote *models.Vote, newVoteType string) error {
	vote.VoteType = newVoteType
	saveErr := database.DB.Save(vote).Error
	if saveErr != nil {
		return errors.New("failed to update vote")
	}
	return nil
}

// GetVoteCountsWithUserVote gets vote counts and user's vote using a single query
func (s *VoteService) GetVoteCountsWithUserVote(votableID, votableType, userID string) (*VoteCounts, error) {
	// where the result is stored
	var result VoteCounts

	// Database query
	err := database.DB.Model(&models.Vote{}).
		// sum case adds up likes and dislikes when vote_type is met and stored as likes and dislikes
		// Coalesce helps to convert any null rows to 0
		// For each row keep only if user id matches, else null and then max just removes all nulls
		Select(`
			COALESCE(SUM(CASE WHEN vote_type = 'like' THEN 1 ELSE 0 END),0) AS likes,
			COALESCE(SUM(CASE WHEN vote_type = 'dislike' THEN 1 ELSE 0 END),0) AS dislikes,
			MAX(CASE WHEN user_id = ? THEN vote_type ELSE NULL END) AS my_vote
		`, userID).
		Where("votable_id = ? AND votable_type = ?", votableID, votableType).
		Scan(&result).Error

	if err != nil {
		return nil, errors.New("failed to get vote counts")
	}

	return &result, nil
}

// GetVoteCounts gets vote counts without user's vote
func (s *VoteService) GetVoteCounts(votableID, votableType string) (int64, int64, error) {
	// variables to hold the amount of likes and dislikes
	var likes int64
	var dislikes int64

	// Count likes from the database
	likeQuery := database.DB.Model(&models.Vote{}).Where("votable_id = ? AND votable_type = ? AND vote_type = ?",
		votableID, votableType, "like").Count(&likes)
	if likeQuery.Error != nil {
		return 0, 0, errors.New("database error")
	}

	// Count dislikes in the database
	dislikeQuery := database.DB.Model(&models.Vote{}).Where("votable_id = ? AND votable_type = ? AND vote_type = ?",
		votableID, votableType, "dislike").Count(&dislikes)
	if dislikeQuery.Error != nil {
		return 0, 0, errors.New("database error")
	}

	return likes, dislikes, nil
}

// ValidateVoteInput validates vote input data
func (s *VoteService) ValidateVoteInput(input VoteInput) error {
	// Validate votable type
	if input.VotableType != "post" && input.VotableType != "comment" {
		return errors.New("invalid content type")
	}

	// Validate vote type
	if input.VoteType != "like" && input.VoteType != "dislike" {
		return errors.New("invalid vote type")
	}

	return nil
}