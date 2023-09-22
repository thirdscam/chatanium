package main

import (
	"antegr.al/chatanium-bot/v1/src/Log"
	"github.com/bwmarrin/discordgo"
)

func main() {
	Log.Init(4)
}

func getClient(token string) *discordgo.Session {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		Log.Error.Fatalln("Failed to create discord session: ", err)
	}

	return discord
}
