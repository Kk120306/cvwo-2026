package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// https://gorm.io/docs/models.html
// Read here for what Gorm Model provides
type Topic struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	Slug      string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// https://gorm.io/docs/hooks.html
// GORM hook that runs before creating Topic - generates a new unique id
func (t *Topic) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New().String()
	return
}
