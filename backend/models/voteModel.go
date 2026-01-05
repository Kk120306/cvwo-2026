package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Vote models represents like or dislike of a post for a user
// This works by storing the userID, the content ID and the type of vote. By doing so, it ensures that a user can strictly have one vote on content
// A vote can then be used to gather the total by finding the count of votes of a certain type on a content
type Vote struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      string    `gorm:"type:uuid;not null;index:unique_vote,unique" json:"userId"`
	VotableID   string    `gorm:"type:uuid;not null;index:unique_vote,unique" json:"votableId"`
	VotableType string    `gorm:"type:varchar(20);not null;index:unique_vote,unique" json:"votableType"`
	VoteType    string    `gorm:"type:varchar(20);not null" json:"voteType"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (v *Vote) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New().String()
	return
}
