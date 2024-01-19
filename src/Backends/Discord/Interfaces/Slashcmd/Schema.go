package Slashcmd

import (
	"github.com/bwmarrin/discordgo"
)

func getDefinedSchema() map[string]func() *discordgo.ApplicationCommand {
	return map[string]func() *discordgo.ApplicationCommand{
		// Write your schema of module here
		"Ping":      GetPing,
		"Snowflake": GetSnowflake,
	}
}

func GetAllSchema() []*discordgo.ApplicationCommand {
	var result []*discordgo.ApplicationCommand
	schemas := getDefinedSchema()

	for _, v := range schemas {
		result = append(result, v())
	}

	return result
}

func GetOnlySchema(AllowedModules []string) []*discordgo.ApplicationCommand {
	var result []*discordgo.ApplicationCommand
	AllSchemas := getDefinedSchema()

	for _, v := range AllowedModules {
		for name, schema := range AllSchemas {
			if name == v {
				result = append(result, schema())
				break
			}
		}
	}

	return result
}
