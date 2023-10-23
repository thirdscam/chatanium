package main

import (
	"flag"

	"antegr.al/chatanium-bot/v1/src/Ignite"
	"antegr.al/chatanium-bot/v1/src/Log"
)

var (
	Token       string
	LoggingMode int
)

func main() {
	flag.StringVar(&Token, "token", "", "Address to proxy")
	flag.IntVar(&LoggingMode, "logging-mode", 4, "Logging mode")

	// Parse the flags and init the logger
	flag.Parse()
	Log.Init(LoggingMode)

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
		Token:    "MTE1NDc4NTkzOTM5OTM3Njk2Ng.GEwjcR.Bc5uPjRJ1ceE8jtkqk3P4iLtCpbPIqx5Gq8brE",
	}
	Discord.Start()
	defer Discord.Shutdown()

	// Wait for a signal to shutdown
	Ignite.WaitSignal()
}
