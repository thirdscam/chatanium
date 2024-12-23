package main

import (
	"os"

	discord "antegr.al/chatanium-bot/v1/src/Backends/Discord"
	Database "antegr.al/chatanium-bot/v1/src/Database"
	util "antegr.al/chatanium-bot/v1/src/Util"
	"antegr.al/chatanium-bot/v1/src/Util/Log"
)

func main() {
	// Init the environment variables
	util.InitEnv()

	// Init the logging system
	Log.Init()

	Log.Info.Println("Antegral/Chatanium: Scalable Bot Management System")
	Log.Info.Println("Press CTRL+C to shutdown.")

	// Ignite Database
	database := Database.DB{}
	database.Start()
	defer database.Shutdown()

	// Ignite Discord
	discord.Start()
	// Discord := Ignite.Discord{
	// 	Database: database.Client,
	// 	Token:    os.Getenv("DISCORD_TOKEN"),
	// }
	// Discord.Start()
	// defer Discord.Shutdown()

	// Wait for a signal to shutdown
	util.WaitSignal()
}
