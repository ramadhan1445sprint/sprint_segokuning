package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig(filename string) {
	godotenv.Load(filename)
}

func GetString(key string) string {
	return os.Getenv(key)
}