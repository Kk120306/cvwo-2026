package database

import (
	"fmt"
	"github.com/Kk120306/cvwo-2026/backend/helpers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

// Using GORM to connect to the database
// Using PostgreSQL as the database - hosted rn on Neon
// https://gorm.io/docs/connecting_to_the_database.html
func ConnectToDb() {
	helpers.LoadEnvVariables()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
