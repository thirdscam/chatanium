package Internal

import (
	"context"
	"time"

	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Log"
	"antegr.al/chatanium-bot/v1/src/util"
	"github.com/bwmarrin/discordgo"
)

func MessageLogger(client *discordgo.Session, database *db.PrismaClient) {
	// TODO: Database (actions)

	// Handle all messages from all guilds
	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == client.State.User.ID {
			return
		}

		createMessage(m, database)
		Log.Verbose.Printf("MESSAGE/SEND (%v -> %v) %v: %v", m.GuildID, m.ChannelID, m.Author.Username, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == client.State.User.ID {
			return
		}

		Log.Verbose.Printf("MESSAGE/UPDATE (%v -> %v / %v) %v: %v", m.GuildID, m.ChannelID, m.Message.ID, m.Author.Username, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageDelete) {
		// FIXME: Cannot get author
		Log.Verbose.Printf("MESSAGE/DELETE (%v -> %v) %v", m.GuildID, m.ChannelID, m.Author)
	})
}

func MemberLogger(client *discordgo.Session, dbClient *db.PrismaClient) {
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

func createMessage(m *discordgo.MessageCreate, database *db.PrismaClient) {
	ctx := context.Background()

	// Database Task: Upsert user
	Users := db.Users
	database.Users.UpsertOne(Users.ID.Equals(util.StringToBigint(m.Author.ID))).Create(
		Users.ID.Set(util.StringToBigint(m.Author.ID)),
		Users.Username.Set(m.Author.Username),
		Users.CreatedAt.Set(time.Now()),
	).Exec(ctx)

	// Database Task: Insert message
	Msg := db.Messages
	_, err := database.Messages.CreateOne(
		Msg.MessageID.Set(util.StringToBigint(m.ID)),
		Msg.Type.Set(int(m.Type)),
		Msg.CreatedAt.Set(m.Timestamp),
		Msg.Users.Link(db.Users.ID.Equals(util.StringToBigint(m.Author.ID))),
		Msg.Guilds.Link(db.Guilds.ID.Equals(util.StringToBigint(m.GuildID))),
		Msg.Channels.Link(db.Channels.ID.Equals(util.StringToBigint(m.ChannelID))),
	).Exec(ctx)
	if err != nil {
		Log.Error.Printf("Failed to create message: %v", err)
	}
}
