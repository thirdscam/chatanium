package Database

import (
	"context"
	"errors"
	"strings"

	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	util "antegr.al/chatanium-bot/v1/src/Util"
	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

func GetMessageInfo(gid, mid string, database *db.PrismaClient) *db.MessagesModel {
	msg, err := database.Messages.FindUnique(
		db.Messages.MessageID.Equals(util.Str2Int64(mid)),
	).With(
		db.Messages.Users.Fetch(),
	).Exec(
		context.Background(),
	)
	if errors.Is(err, db.ErrNotFound) {
		Log.Error.Printf("G:%s | M:%s > Cannot find message : %v", gid, mid, err)
		return nil
	} else if err != nil {
		Log.Error.Printf("G:%s | M:%s > Failed to find message : %v", gid, mid, err)
		return nil
	}

	return msg
}

func CreateMessage(s *discordgo.Session, m *discordgo.MessageCreate, database *db.PrismaClient) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	ctx := context.Background()

	// Database insert user
	InsertUser(database, m.Author.ID, m.Author.Username)

	// Database insert member (guild user)
	InsertMember(database, m.Author.ID, m.GuildID, m.Message.Member.Nick)

	// Database Task: Insert message
	Msg := db.Messages
	_, err := database.Messages.CreateOne(
		Msg.MessageID.Set(util.Str2Int64(m.ID)),
		Msg.Type.Set(int(m.Type)),
		Msg.CreatedAt.Set(m.Timestamp),
		Msg.Users.Link(db.Users.ID.Equals(util.Str2Int64(m.Author.ID))),
		Msg.Contents.Set(m.Content),
		Msg.Guilds.Link(db.Guilds.ID.Equals(util.Str2Int64(m.GuildID))),
		Msg.Channels.Link(db.Channels.ID.Equals(util.Str2Int64(m.ChannelID))),
	).Exec(ctx)
	if err != nil {
		Log.Error.Printf("Failed to create message: %v", err)
	}
}

func UpdateMessage(s *discordgo.Session, m *discordgo.MessageUpdate, database *db.PrismaClient) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
}

func DeleteMessage(s *discordgo.Session, m *discordgo.MessageDelete, database *db.PrismaClient) {
	msg := GetMessageInfo(m.GuildID, m.ID, database)
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
