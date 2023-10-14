package Handlers

import "github.com/bwmarrin/discordgo"

func GetAll() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		// Write your handler of COMMANDS for module here
		"ping":           HandlePing,
		"snowflake2time": HandleSnowflake2time,
	}
}

func GetAllowedOnly(AllowedModules []string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	result := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	AllHandlers := GetAll()

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
