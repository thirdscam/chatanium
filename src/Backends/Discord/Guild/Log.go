package Guild

import (
	backendDB "antegr.al/chatanium-bot/v1/src/Backends/Discord/Database"
	"antegr.al/chatanium-bot/v1/src/Database"
	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/bwmarrin/discordgo"
)

func HandleMessageLog(client *discordgo.Session, db *Database.DB) {
	// Handle all messages from all guilds
	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		backendDB.CreateMessage(s, m)
		Log.Verbose.Printf("G:%v | C:%v > %v: %v", m.GuildID, m.ChannelID, m.Author.Username, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		backendDB.UpdateMessage(s, m, database)
		Log.Verbose.Printf("G:%v | C:%v > Update M:%v > %v", m.GuildID, m.ChannelID, m.Message.ID, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageDelete) {
		backendDB.DeleteMessage(s, m, database)
		Log.Verbose.Printf("G:%v | C:%v > Delete M:%v", m.GuildID, m.ChannelID, m.ID)
	})
}

func HandleMemberLog(client *discordgo.Session, db *Database.DB) {
	// TODO: Database (Save Message, actions)
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
