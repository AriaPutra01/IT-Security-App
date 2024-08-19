package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariable() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
