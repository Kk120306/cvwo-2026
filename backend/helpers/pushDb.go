package helpers

import (
	"github.com/Kk120306/cvwo-2026/backend/models"
)

// https://gorm.io/docs/index.html -- Migrating the schema
func PushDb() {
	DB.AutoMigrate(&models.User{})
}
