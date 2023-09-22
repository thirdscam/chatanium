package Guild

import (
	"antegr.al/chatanium-bot/v1/src/Handlers"
	"antegr.al/chatanium-bot/v1/src/Log"
	"antegr.al/chatanium-bot/v1/src/Schema"
	"github.com/bwmarrin/discordgo"
)

func Handle(client *discordgo.Session) {
	// TODO: Database (Save Guild ID, etc.)
	var GuildCmds []Commands

	// Register all commands from all guilds
	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		Log.Verbose.Printf("Joined Guild: %v (%v)", g.Name, g.ID)

		AllowedModules := []string{"ping"}

		Guild := Commands{
			Schema:   Schema.GetAllowedOnly(AllowedModules),
			Handlers: Handlers.GetAllowedOnly(AllowedModules),
			Client:   client,
			GuildID:  g.ID,
		}

		// Register commands from guild
		Guild.RegisterHandlers()
		Guild.RegisterSchema()

		GuildCmds = append(GuildCmds, Guild)
	})

	// Handle all messages from all guilds
	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == client.State.User.ID {
			return
		}

		Log.Verbose.Printf("(%v -> %v) %v: %v", m.GuildID, m.ChannelID, m.Author, m.Content)
	})
}
