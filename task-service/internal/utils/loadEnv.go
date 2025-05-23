package utils

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

// return env file path
func LoadEnv() (string, error) {
	envFilePath := ".env"
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatalf("Error loading .env file")
		return "", err
	}
	fmt.Println("Successfully loaded .env file")
	return envFilePath, nil
}
