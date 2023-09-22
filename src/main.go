package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"antegr.al/chatanium-bot/v1/src/Guild"
	"antegr.al/chatanium-bot/v1/src/Ignite"
	"antegr.al/chatanium-bot/v1/src/Log"
	"github.com/bwmarrin/discordgo"
)

var (
	Token       string
	LoggingMode int
)

func main() {
	flag.StringVar(&Token, "token", "", "Address to proxy")
	flag.IntVar(&LoggingMode, "logging-mode", 3, "Logging mode")

	// Parse the flags and init the logger
	flag.Parse()
	Log.Init(LoggingMode)

	// Create the session and the registered commands list
	client := getClient("MTE1NDc4NTkzOTM5OTM3Njk2Ng.GEwjcR.Bc5uPjRJ1ceE8jtkqk3P4iLtCpbPIqx5Gq8brE")
	var RegisteredGuildCmds *[]Guild.Commands

	// Create a channel to receive OS signals
	stop := make(chan os.Signal)

	// Ignite Server (Discord Bot, Status Page, etc.)
	go Ignite.Discord(stop, client, RegisteredGuildCmds)

	// Wait for a signal to shutdown
	log.Println("Press Ctrl+C to shutdown")
	<-stop
	shutdown(stop, client, RegisteredGuildCmds)
}

func shutdown(Singal chan os.Signal, Client *discordgo.Session, RegisteredGuildCmds *[]Guild.Commands) {
	Log.Info.Println("Shutting down...")

	if RegisteredGuildCmds != nil {
		// Remove all commands from all guilds
		for _, v := range *RegisteredGuildCmds {
			v.RemoveSchema()
		}
	}

	// Close the client
	Client.Close()

	Log.Info.Println("Successfully shutdown.")
	signal.Notify(Singal, os.Interrupt)
}

func getClient(token string) *discordgo.Session {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		Log.Error.Fatalln("Failed to create discord session: ", err)
	}

	return discord
}
