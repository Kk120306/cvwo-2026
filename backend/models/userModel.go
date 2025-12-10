package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// https://gorm.io/docs/models.html
// Read here for what Gorm Model provides
type User struct {
	ID                string `gorm:"type:uuid;primaryKey"`
	Username          string `gorm:"uniqueIndex;not null"`
	AvatarURL         string `gorm:"default:'https://localhost:3000/placeholder.png'"`
	IsAdmin           bool   `gorm:"default:false"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// https://gorm.io/docs/hooks.html
// GORM hook that runs before creating user - generates a new unique id
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
