package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func Init() {
	if os.Getenv("GO_ENV") == "development" || os.Getenv("GO_ENV") == "production" {
		log.Println("Now you use docker-compose environment")
	} else {
		ex, _ := os.Executable()
		exPath := filepath.Dir(ex)
		if err := godotenv.Load(exPath + "/.env"); err != nil {
			log.Fatal(".env is not loaded properly\n", err.Error())
		}
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
