package Command

import (
	"antegr.al/chatanium-bot/v1/src/Log"
	"github.com/bwmarrin/discordgo"
)

type GuildCommands struct {
	Schema   []*discordgo.ApplicationCommand
	Handlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	Client   *discordgo.Session
	GuildID  string
}

func (t *GuildCommands) RegisterHandlers() {
	t.Client.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := t.Handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func (t *GuildCommands) RegisterSchema(Commands []*discordgo.ApplicationCommand) {
	Log.Info.Printf("%s: Adding commands...", t.GuildID)
	registered := make([]*discordgo.ApplicationCommand, len(t.Schema))

	for i, v := range t.Schema {
		cmd, err := t.Client.ApplicationCommandCreate(t.Client.State.User.ID, t.GuildID, v)
		if err != nil {
			Log.Error.Fatalf("%s: Cannot create '%v' command: %v", t.GuildID, v.Name, err)
		}
		registered[i] = cmd
	}
}

func (t *GuildCommands) RemoveSchema() {
	Log.Info.Printf("%s Removing commands...", t.GuildID)

	for _, v := range t.Schema {
		err := t.Client.ApplicationCommandDelete(t.Client.State.User.ID, t.GuildID, v.ID)
		if err != nil {
			Log.Error.Fatalf("%s: Cannot delete '%v' command: %v", t.GuildID, v.Name, err)
		}
	}
}
