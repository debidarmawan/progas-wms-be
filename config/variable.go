package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Init() {
	// Load .env for local dev. Variables already set in the OS (Docker/K8s) are not overridden.
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		log.Fatal("failed to load .env: ", err.Error())
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
