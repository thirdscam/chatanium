package Command

import (
	"log"

	"antegr.al/chatanium-bot/v1/src/Log"
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Commands []*discordgo.ApplicationCommand
	Client   *discordgo.Session
	GuildID  string
}

func (t *Command) SetAll(Commands []*discordgo.ApplicationCommand) {
	Log.Info.Println("Adding commands...")
	registered := make([]*discordgo.ApplicationCommand, len(t.Commands))

	for i, v := range t.Commands {
		cmd, err := t.Client.ApplicationCommandCreate(t.Client.State.User.ID, t.GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registered[i] = cmd
	}
}

func (t *Command) RemoveAll(RegisteredCommands []*discordgo.ApplicationCommand) {
	log.Println("Removing commands...")

	for _, v := range RegisteredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
