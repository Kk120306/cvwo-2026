package database

import (
	"github.com/Kk120306/cvwo-2026/backend/models"
)

// https://gorm.io/docs/index.html -- Migrating the schema
func PushDb() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Topic{})
	DB.AutoMigrate(&models.Post{})
	DB.AutoMigrate(&models.Comment{})
	DB.AutoMigrate(&models.Vote{})
}
