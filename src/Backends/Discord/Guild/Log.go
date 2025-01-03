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
		backendDB.CreateMessage(s, m, db.Queries)
		Log.Verbose.Printf("G:%v | C:%v > %v: %v", m.GuildID, m.ChannelID, m.Author.Username, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		backendDB.UpdateMessage(s, m, db.Queries)
		Log.Verbose.Printf("G:%v | C:%v > UPDATE > M:%v (=> %v)", m.GuildID, m.ChannelID, m.Message.ID, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageDelete) {
		backendDB.DeleteMessage(s, m, db.Queries)
		Log.Verbose.Printf("G:%v | C:%v > UPDATE > M:%v", m.GuildID, m.ChannelID, m.ID)
	})
}

func HandleMemberLog(client *discordgo.Session, db *Database.DB) {
	// TODO: Database (Save Message, actions)
	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		Log.Verbose.Printf("G:%v > MEMBER_ADD > U:%v (%v)", m.GuildID, m.User.ID, m.Nick)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
		Log.Verbose.Printf("G:%v > MEMBER_REMOVE > U:%v (%v)", m.GuildID, m.User.ID, m.Nick)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
		Log.Verbose.Printf("G:%v > MEMBER_UPDATE > U:%v (%v)", m.GuildID, m.User.ID, m.Nick)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildBanAdd) {
		Log.Warn.Printf("G:%v > BAN_ADD > U:%v", m.GuildID, m.User.ID)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.GuildBanRemove) {
		Log.Warn.Printf("G:%v > BAN_REMOVE > U:%v", m.GuildID, m.User.ID)
	})
}
