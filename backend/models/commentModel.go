package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Consider adding nesting later on
type Comment struct {
	ID       string `gorm:"type:uuid;primaryKey" json:"id"`
	PostID   string `gorm:"type:uuid;not null" json:"postId"`
	AuthorID string `gorm:"type:uuid;not null" json:"authorId"`

	Post   Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Author User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`

	Content  string `gorm:"type:text;not null" json:"content"`
	IsPinned bool   `gorm:"default:false" json:"isPinned"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Deletes any related field with cascade
	Votes []Vote `gorm:"foreignKey:VotableID;constraint:OnDelete:CASCADE" json:"votes,omitempty"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}
