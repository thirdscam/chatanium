package Util

import (
	"os"

	"antegr.al/chatanium-bot/v1/src/Log"
)

// Initalize the environment variables from .env file
func InitLog() {
	mode := os.Getenv("LOG_MODE")
	switch mode {
	case "production":
		Log.Init(Log.PRODUCTION_MODE)
	case "development":
		Log.Init(Log.DEVELOPMENT_MODE)
	default:
		Log.Init(Log.PRODUCTION_MODE)
		Log.Warn.Printf("Invalid logging mode: %s", mode)
	}
}
