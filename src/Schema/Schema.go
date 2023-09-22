package Commands

import "github.com/bwmarrin/discordgo"

func GetPing() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "ping to bot.",
	}
}
