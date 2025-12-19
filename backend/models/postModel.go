package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// https://gorm.io/docs/models.html
// Read here for what Gorm Model provides
type Post struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	TopicID   string    `gorm:"type:uuid;not null"`
	AuthorID  string    `gorm:"type:uuid;not null"`
	Topic     Topic     `gorm:"foreignKey:TopicID"`
	Author    User      `gorm:"foreignKey:AuthorID"`
	Title     string    `gorm:"type:varchar(255);not null"`
	Content   string    `gorm:"type:text;not null"`
	IsPinned  bool      `gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ImageUrl  *string    `gorm:"type:text"`

	// Deletes any related field with cascade
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Votes    []Vote    `gorm:"foreignKey:VotableID;constraint:OnDelete:CASCADE"`
}

// https://gorm.io/docs/hooks.html
// GORM hook that runs before creating Topic - generates a new unique id
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New().String()
	return
}
