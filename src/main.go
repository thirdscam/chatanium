package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"antegr.al/chatanium-bot/v1/src/Guild"
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

	client := getClient("MTE1NDc4NTkzOTM5OTM3Njk2Ng.GEwjcR.Bc5uPjRJ1ceE8jtkqk3P4iLtCpbPIqx5Gq8brE")
	var RegisteredGuildCmds *[]Guild.Commands

	stop := make(chan os.Signal)
	go start(stop, client, RegisteredGuildCmds)
	go shutdown(stop, client, RegisteredGuildCmds)

	log.Println("Press Ctrl+C to shutdown")
	<-stop
}

func start(singal chan os.Signal, client *discordgo.Session, RegisteredGuildCmds *[]Guild.Commands) {
	Log.Info.Println("Starting Bot...")

	singal <- os.Interrupt
}

func shutdown(Singal chan os.Signal, client *discordgo.Session, RegisteredGuildCmds *[]Guild.Commands) {
	Log.Info.Println("Shutting down...")

	// Remove all commands from all guilds
	for _, v := range *RegisteredGuildCmds {
		v.RemoveSchema()
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
