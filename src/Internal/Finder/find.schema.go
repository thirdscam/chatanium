package Schema

import "github.com/bwmarrin/discordgo"

func GetFind() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "find",
		Description: "Find a user, channel, or guild by ID",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "Enter a ID (C:1154789841654009886, M:1171064681746681926, etc.)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "expose",
				Description: "Expose the result to everyone, default is false",
				Required:    false,
			},
		},
	}
}
