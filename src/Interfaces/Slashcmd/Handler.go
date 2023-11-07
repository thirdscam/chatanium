package Slashcmd

import (
	"github.com/bwmarrin/discordgo"
)

func GetAllHandler() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		// Write your handler of COMMANDS for module here
		"ping":           HandlePing,
		"snowflake2time": HandleSnowflake2time,
	}
}

func GetOnlyHandler(AllowedModules []string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	result := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	AllHandlers := GetAllHandler()

	for _, v := range AllowedModules {
		for name, fn := range AllHandlers {
			if name == v {
				result[v] = fn
				break
			}
		}
	}

	return result
}
