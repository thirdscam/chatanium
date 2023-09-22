package Handlers

import (
	"github.com/bwmarrin/discordgo"
)

func GetHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": HandlePing,
	}
}
