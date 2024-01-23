package Guild

import (
	"antegr.al/chatanium-bot/v1/src/Backends/Discord/Database"
	slash "antegr.al/chatanium-bot/v1/src/Backends/Discord/Interfaces/Slash"
	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/bwmarrin/discordgo"
)

// Handle all events from guild.
// This function is used for entry point of discord backend.
func Handle(client *discordgo.Session, db *db.PrismaClient) {
	/******************** Interfaces ********************/
	slash := slash.Guild{
		Client: client,
	}
	slash.Start() // Start slash command manager

	/******************** Guild Events ********************/
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

	/******************** Chat Events ********************/
	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		createMessage(s, m, db)
		Log.Verbose.Printf("G:%v | C:%v > %v: %v", m.GuildID, m.ChannelID, m.Author.Username, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		updateMessage(s, m, db)
		Log.Verbose.Printf("G:%v | C:%v > Update M:%v > %v", m.GuildID, m.ChannelID, m.Message.ID, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageDelete) {
		deleteMessage(s, m, db)
		Log.Verbose.Printf("G:%v | C:%v > Delete M:%v", m.GuildID, m.ChannelID, m.ID)
	})

	/******************** Member Events ********************/
	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		Log.Verbose.Printf("MEMBER/JOIN (%v) %v", m.GuildID, m.Nick)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
		Log.Verbose.Printf("MEMBER/LEFT (%v) %v", m.GuildID, m.Nick)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
		Log.Verbose.Printf("MEMBER/UPDATE (%v) %v", m.GuildID, m.Nick)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildBanAdd) {
		Log.Warn.Printf("MEMBER/BAN (%v) %v(%v)", m.GuildID, m.User.Username, m.User.ID)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildBanRemove) {
		Log.Warn.Printf("MEMBER/UNBAN (%v) %v(%v)", m.GuildID, m.User.Username, m.User.ID)
	})
}
