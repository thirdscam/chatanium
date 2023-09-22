package Schema

import "github.com/bwmarrin/discordgo"

func GetAll() []*discordgo.ApplicationCommand {
	modules := map[string]func() *discordgo.ApplicationCommand{
		// Write your schema of commands for module here
		"ping": GetPing,
	}

	var result []*discordgo.ApplicationCommand

	for _, v := range modules {
		result = append(result, v())
	}

	return result
}
