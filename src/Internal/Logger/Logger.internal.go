package Logger

import (
	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Log"
	"github.com/bwmarrin/discordgo"
)

func Message(client *discordgo.Session, database *db.PrismaClient) {
	// TODO: Database (actions)

	// Handle all messages from all guilds
	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		createMessage(s, m, database)
		Log.Verbose.Printf("G:%v | C:%v > %v: %v", m.GuildID, m.ChannelID, m.Author.Username, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		updateMessage(s, m, database)
		Log.Verbose.Printf("G:%v | C:%v > Update M:%v > %v", m.GuildID, m.ChannelID, m.Message.ID, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageDelete) {
		deleteMessage(s, m, database)
		Log.Verbose.Printf("G:%v | C:%v > Delete M:%v", m.GuildID, m.ChannelID, m.ID)
	})
}

func Member(client *discordgo.Session, dbClient *db.PrismaClient) {
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
