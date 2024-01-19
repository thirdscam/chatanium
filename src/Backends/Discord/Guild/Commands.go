package Guild

import (
	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/bwmarrin/discordgo"
)

type Commands struct {
	Schema   []*discordgo.ApplicationCommand
	Handlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	Client   *discordgo.Session
	GuildID  string
}

func (t *Commands) RegisterHandlers() {
	Log.Verbose.Printf("%s > Adding handlers...", t.GuildID)
	t.Client.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := t.Handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func (t *Commands) RegisterSchema() {
	Log.Verbose.Printf("%s > Adding commands...", t.GuildID)
	registered := make([]*discordgo.ApplicationCommand, len(t.Schema))

	for i, v := range t.Schema {
		cmd, err := t.Client.ApplicationCommandCreate(t.Client.State.User.ID, t.GuildID, v)
		if err != nil {
			Log.Error.Fatalf("%s: Cannot create '%v' command: %v", t.GuildID, v.Name, err)
		}
		registered[i] = cmd
	}
}

func (t *Commands) RemoveSchema() {
	Log.Verbose.Printf("%s > Removing commands...", t.GuildID)

	for _, v := range t.Schema {
		err := t.Client.ApplicationCommandDelete(t.Client.State.User.ID, t.GuildID, v.ID)
		if err != nil {
			Log.Error.Fatalf("%s: Cannot delete '%v' command: %v", t.GuildID, v.Name, err)
		}
	}
}
