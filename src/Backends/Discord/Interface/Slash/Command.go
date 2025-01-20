package Slash

import (
	"github.com/bwmarrin/discordgo"
)

type Commands map[*discordgo.ApplicationCommand]func(s *discordgo.Session, i *discordgo.InteractionCreate)

var CommandMap Commands

// get commands for guild from pre-defined commands and modules.
func getCommands(guildId string) Commands {
	// TODO(Security): Support ACL
	_ = guildId
	return CommandMap
}
