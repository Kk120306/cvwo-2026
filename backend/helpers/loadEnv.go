package helpers

import (
	"github.com/joho/godotenv"
	"log"
)

// Loads enviornement variables - public function
// https://github.com/joho/godotenv - just taken from here
func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
