package Logger

import (
	"context"
	"strings"

	"antegr.al/chatanium-bot/v1/src/Database"
	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Log"
	util "antegr.al/chatanium-bot/v1/src/Util"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

func createMessage(s *discordgo.Session, m *discordgo.MessageCreate, database *db.PrismaClient) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	ctx := context.Background()

	// Database insert user
	Database.InsertUser(database, m.Author.ID, m.Author.Username)

	// Database insert member (guild user)
	Database.InsertMember(database, m.Author.ID, m.GuildID, m.Message.Member.Nick)

	// Database Task: Insert message
	Msg := db.Messages
	_, err := database.Messages.CreateOne(
		Msg.MessageID.Set(util.StringToBigint(m.ID)),
		Msg.Type.Set(int(m.Type)),
		Msg.CreatedAt.Set(m.Timestamp),
		Msg.Users.Link(db.Users.ID.Equals(util.StringToBigint(m.Author.ID))),
		Msg.Contents.Set(m.Content),
		Msg.Guilds.Link(db.Guilds.ID.Equals(util.StringToBigint(m.GuildID))),
		Msg.Channels.Link(db.Channels.ID.Equals(util.StringToBigint(m.ChannelID))),
	).Exec(ctx)
	if err != nil {
		Log.Error.Printf("Failed to create message: %v", err)
	}
}

func updateMessage(s *discordgo.Session, m *discordgo.MessageUpdate, database *db.PrismaClient) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
}

func deleteMessage(s *discordgo.Session, m *discordgo.MessageDelete, database *db.PrismaClient) {
	msg := Database.GetMessageInfo(m.GuildID, m.ID, database)
	if msg == nil {
		Log.Warn.Println("MessageIntegrityCheck: Cannot found message. may created when bot was offline.")
		Log.Warn.Printf("G:%v | C:%v > Cannot found message at M:%v", m.GuildID, m.ChannelID, m.ID)
		return
	}

	content, err := msg.Contents()
	if err != true {
		Log.Error.Printf("G:%v | C:%v > Failed to get message content from M:%v", m.GuildID, m.ChannelID, err)
		return
	}

	Log.Info.Printf(color.RedString("G:%v | C:%v > ACTION/DELETE > %v: %v", m.GuildID, m.ChannelID, strings.TrimRight(msg.Users().Username, " "), content))
}
