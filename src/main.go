package main

import (
	"os"

	"antegr.al/chatanium-bot/v1/src/Ignite"
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

	// Example: Ignite Backend for 3 steps
	// Database := Ignite.DB{}   // 1. Prepare
	// db := Database.Start()    // 2. Start
	// defer Database.Shutdown() // 3. Shutdown

	// Ignite Database
	Database := Ignite.DB{}
	db := Database.Start()
	defer Database.Shutdown()

	// Ignite Discord
	Discord := Ignite.Discord{
		Database: db,
		Token:    os.Getenv("DISCORD_TOKEN"),
	}
	Discord.Start()
	defer Discord.Shutdown()

	// Wait for a signal to shutdown
	Ignite.WaitSignal()
}
