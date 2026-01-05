package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// https://gorm.io/docs/models.html
// Read here for what Gorm Model provides
type Post struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	TopicID   string    `gorm:"type:uuid;not null" json:"topicId"`
	AuthorID  string    `gorm:"type:uuid;not null" json:"authorId"`
	Topic     Topic     `gorm:"foreignKey:TopicID" json:"topic,omitempty"`
	Author    User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	IsPinned  bool      `gorm:"default:false" json:"isPinned"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ImageUrl  *string   `gorm:"type:text" json:"imageUrl,omitempty"`

	// Deletes any related field with cascade
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
	Votes    []Vote    `gorm:"foreignKey:VotableID;constraint:OnDelete:CASCADE" json:"votes,omitempty"`
}

// https://gorm.io/docs/hooks.html
// GORM hook that runs before creating Topic - generates a new unique id
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New().String()
	return
}
