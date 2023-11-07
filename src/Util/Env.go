package Util

import (
	"os"

	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/joho/godotenv"
)

func InitEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		Log.Error.Fatalf("Error loading .env file: %v", err)
	}
}

func GetEnv(EnvName string) string {
	env := os.Getenv(EnvName)
	if env == "" {
		Log.Warn.Printf("CheckEnvExists: Unable to find environment variable: %s", EnvName)
	}
	return env
}
