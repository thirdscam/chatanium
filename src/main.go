package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

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

	// client := getClient("MTE1NDc4NTkzOTM5OTM3Njk2Ng.GEwjcR.Bc5uPjRJ1ceE8jtkqk3P4iLtCpbPIqx5Gq8brE")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Println("Press Ctrl+C to exit")
	<-stop
}

func getClient(token string) *discordgo.Session {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		Log.Error.Fatalln("Failed to create discord session: ", err)
	}

	return discord
}
