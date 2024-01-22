package Guild

import (
	"antegr.al/chatanium-bot/v1/src/Backends/Discord/Database"
	slash "antegr.al/chatanium-bot/v1/src/Backends/Discord/Interfaces/Slash"
	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/bwmarrin/discordgo"
)

func Handle(client *discordgo.Session, db *db.PrismaClient) {
	slash := slash.Guild{
		Client: client,
	}

	// Guild Events
	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		Log.Verbose.Printf("Join Guild: %v (%v)", g.Name, g.ID)
		Database.RegisterGuild(client, db, g.ID, g.OwnerID) // Register guild to database
		slash.OnGuildCreated(g.ID)                          // Register slash commands
	})

	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildDelete) {
		Log.Verbose.Printf("Left Guild: %v (%v)", g.Name, g.ID)
		slash.OnGuildDeleted(g.ID) // Remove slash commands
	})

	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildUpdate) {
		Log.Verbose.Printf("Updated Guild: %v (%v)", g.Name, g.ID)
		// Database.UpdateGuild(client, db, g.ID, g.OwnerID) // TODO(Feature): Update guild to database
		// slash.OnGuildUpdated(g.ID) // TODO(Feature): Update slash commands
	})

	// Internal Modules
	// TODO(Feature): Support internal modules
	// Logger.Member(client, db)
	// Logger.Message(client, db)
}
