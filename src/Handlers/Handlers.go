package Handlers

import "github.com/bwmarrin/discordgo"

func GetAll() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		// Write your handler of commands for module here
		"ping": HandlePing,
	}
}
