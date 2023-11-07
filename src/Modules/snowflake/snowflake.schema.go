package Snowflake

import "github.com/bwmarrin/discordgo"

func GetSnowflake() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "snowflake2time",
		Description: "Discord Snowflake ID to Timestamp",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "snowflake",
				Description: "Enter a Snowflake ID",
				Required:    true,
			},
		},
	}
}
