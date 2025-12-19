package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Consider adding nesting later on
type Comment struct {
	ID       string `gorm:"type:uuid;primaryKey"`
	PostID   string `gorm:"type:uuid;not null"`
	AuthorID string `gorm:"type:uuid;not null"`

	Post   Post `gorm:"foreignKey:PostID"`
	Author User `gorm:"foreignKey:AuthorID"`

	Content  string `gorm:"type:text;not null"`
	IsPinned bool   `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Deletes any related field with cascade
	Votes []Vote `gorm:"foreignKey:VotableID;constraint:OnDelete:CASCADE"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}
