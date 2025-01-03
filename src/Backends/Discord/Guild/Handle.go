package Guild

import (
	backendDB "antegr.al/chatanium-bot/v1/src/Backends/Discord/Database"
	"antegr.al/chatanium-bot/v1/src/Database"
	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/bwmarrin/discordgo"
)

// Handle all events from guild.
// This function is used for entry point of discord backend.
func Handle(client *discordgo.Session, db *Database.DB) {
	/******************** Interfaces ********************/
	// Slash := slash.Guild{
	// 	Client: client,
	// }
	// Slash.Start() // Start slash command manager

	/******************** Guild Events ********************/
	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		Log.Verbose.Printf("G:%v > A:GUILD_CREATE > %v", g.ID, g.Name)
		backendDB.RegisterGuild(client, db.Conn, db.Queries, g.ID, g.OwnerID) // Register guild to database
		// Slash.OnGuildCreated(g.ID)                          // Register slash commands
	})

	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildDelete) {
		Log.Verbose.Printf("G:%v > A:GUILD_DELETE > %v", g.ID, g.Name)
		// Slash.OnGuildDeleted(g.ID) // Remove slash commands
	})

	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildUpdate) {
		Log.Verbose.Printf("G:%v > A:GUILD_UPDATE > %v", g.ID, g.Name)
		// Database.UpdateGuild(client, db, g.ID, g.OwnerID) // TODO(Feature): Update guild to database
		// slash.OnGuildUpdated(g.ID) // TODO(Feature): Update slash commands
	})

	/******************** Chat Events ********************/
	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		backendDB.CreateMessage(s, m, db.Queries)
		Log.Verbose.Printf("G:%v | C:%v > A:CHAT_CREATE > %v: %v", m.GuildID, m.ChannelID, m.Author.Username, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		backendDB.UpdateMessage(s, m, db.Queries)
		Log.Verbose.Printf("G:%v | C:%v > A:CHAT_UPDATE > M:%v (=> %v)", m.GuildID, m.ChannelID, m.Message.ID, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageDelete) {
		backendDB.DeleteMessage(s, m, db.Queries)
		Log.Verbose.Printf("G:%v | C:%v > A:CHAT_DELETE > M:%v", m.GuildID, m.ChannelID, m.ID)
	})

	/******************** Member Events ********************/
	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		Log.Verbose.Printf("G:%v > A:MEMBER_ADD > U:%v (%v)", m.GuildID, m.User.ID, m.Nick)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
		Log.Verbose.Printf("G:%v > A:MEMBER_REMOVE > U:%v (%v)", m.GuildID, m.User.ID, m.Nick)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
		Log.Verbose.Printf("G:%v > A:MEMBER_UPDATE > U:%v (%v)", m.GuildID, m.User.ID, m.Nick)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildBanAdd) {
		Log.Warn.Printf("G:%v > A:BAN_ADD > U:%v", m.GuildID, m.User.ID)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildBanRemove) {
		Log.Warn.Printf("G:%v > A:BAN_REMOVE > U:%v", m.GuildID, m.User.ID)
	})
}
