package Internal

import (
	"antegr.al/chatanium-bot/v1/src/Log"
	"github.com/bwmarrin/discordgo"
)

func MessageLogger(client *discordgo.Session) {
	// TODO: Database (Save Meesage, actions)

	// Handle all messages from all guilds
	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == client.State.User.ID {
			return
		}

		Log.Verbose.Printf("(%v -> %v) %v: %v", m.GuildID, m.ChannelID, m.Author, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == client.State.User.ID {
			return
		}

		Log.Verbose.Printf("MESSAGE/UPDATE (%v -> %v) %v: %v", m.GuildID, m.ChannelID, m.Author, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageDelete) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == client.State.User.ID {
			return
		}

		Log.Verbose.Printf("MESSAGE/DELETE (%v -> %v) %v: %v", m.GuildID, m.ChannelID, m.Author, m.Content)
	})
}

func MemberLogger(client *discordgo.Session) {
	// TODO: Database (Save Meesage, actions)

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
