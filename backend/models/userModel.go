package models

import (
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
// Read here for what Gorm Model provides
type User struct {
	gorm.Model
	Email             string `gorm:"uniqueIndex"`
	EncryptedPassword string `gorm:"not null"` // bcrypt hashed
}
