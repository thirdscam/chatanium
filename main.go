package main

import (
	_ "github.com/joho/godotenv/autoload"
	discord "github.com/thirdscam/chatanium/src/Backends/Discord"
	Database "github.com/thirdscam/chatanium/src/Database"
	util "github.com/thirdscam/chatanium/src/Util"
	"github.com/thirdscam/chatanium/src/Util/Log"
)

func main() {
	// Init the logging system
	Log.Init()

	// Init the environment variables
	util.InitEnv()

	Log.Info.Println("Antegral/Chatanium: Scalable Bot Management System")
	Log.Info.Println("Press CTRL+C to shutdown.")

	// Ignite Database
	database := &Database.DB{}
	database.Start()
	defer database.Shutdown()

	// Create a new Discord session using the provided bot token.
	discord := &discord.Backend{
		Token:    util.GetEnv("DISCORD_TOKEN"),
		Database: database,
	}
	discord.Start()
	defer discord.Shutdown()

	// Wait for a signal to shutdown
	util.WaitSignal()
}
