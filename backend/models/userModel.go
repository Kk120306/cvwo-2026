package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// https://gorm.io/docs/models.html
// Read here for what Gorm Model provides
type User struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	AvatarURL string    `gorm:"default:'https://d1nxlczpemry9k.cloudfront.net/829472_man_512x512.png'" json:"avatarUrl"`
	IsAdmin   bool      `gorm:"default:false" json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// https://gorm.io/docs/hooks.html
// GORM hook that runs before creating user - generates a new unique id
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
