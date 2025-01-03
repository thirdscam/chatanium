package Database

import (
	"context"
	"database/sql"
	"errors"

	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	util "antegr.al/chatanium-bot/v1/src/Util"
	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

func GetMessageInfo(queries *db.Queries, mid string) *db.Message {
	msg, err := queries.GetMessage(context.Background(), util.Str2Int64(mid))
	if errors.Is(err, sql.ErrNoRows) {
		Log.Error.Printf("M:%s > Cannot find message : %v", mid, err)
		return nil
	} else if err != nil {
		Log.Error.Printf("M:%s > Failed to find message: %v", mid, err)
		return nil
	}

	return &msg
}

func CreateMessage(s *discordgo.Session, m *discordgo.MessageCreate, queries *db.Queries) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Database insert user
	InsertUser(queries, m.Author.ID, m.Author.Username)

	// Database insert member (guild user)
	InsertMember(queries, m.Author.ID, m.GuildID, m.Message.Member.Nick)

	// Database Task: Insert message
	if err := queries.InsertMessage(context.Background(), db.InsertMessageParams{
		MessageID: util.Str2Int64(m.ID),
		Type:      int64(m.Type),
		CreatedAt: m.Timestamp,
		UserID:    util.Str2Int64(m.Author.ID),
		Contents: sql.NullString{
			String: m.Content,
			Valid:  true,
		},
		GuildID: sql.NullInt64{
			Int64: util.Str2Int64(m.GuildID),
			Valid: true,
		},
		ChannelID: sql.NullInt64{
			Int64: util.Str2Int64(m.ChannelID),
			Valid: true,
		},
	}); err != nil {
		Log.Error.Printf("M:%s > Failed to create message: %v", m.ID, err)
	}
}

func UpdateMessage(s *discordgo.Session, m *discordgo.MessageUpdate, queries *db.Queries) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// if err := queries.InsertMessage(context.Background(), db.InsertMessageParams{
	// 	MessageID: util.Str2Int64(m.ID),
	// 	Type:      int64(m.Type),
	// 	CreatedAt: m.Timestamp,
	// 	UserID:    util.Str2Int64(m.Author.ID),
	// 	Contents: sql.NullString{
	// 		String: m.Content,
	// 		Valid:  true,
	// 	},
	// 	GuildID: sql.NullInt64{
	// 		Int64: util.Str2Int64(m.GuildID),
	// 		Valid: true,
	// 	},
	// 	ChannelID: sql.NullInt64{
	// 		Int64: util.Str2Int64(m.ChannelID),
	// 		Valid: true,
	// 	},
	// }); err != nil {
	// 	Log.Error.Printf("M:%s > Failed to update message: %v", m.ID, err)
	// }
}

func DeleteMessage(s *discordgo.Session, m *discordgo.MessageDelete, queries *db.Queries) {
	msg, err := queries.GetMessage(context.Background(), util.Str2Int64(m.ID))

	if errors.Is(err, sql.ErrNoRows) {
		Log.Warn.Println("MessageIntegrityCheck: Cannot find message. may created when bot was offline.")
		Log.Warn.Printf("G:%v | C:%v > Cannot find message at M:%v", m.GuildID, m.ChannelID, m.ID)
	} else if err != nil {
		Log.Error.Printf("M:%s > Failed to find message: %v", m.ID, err)
	}

	if err := queries.DeleteMessage(context.Background(), util.Str2Int64(m.ID)); err != nil {
		Log.Error.Printf("M:%s > Failed to delete message: %v", m.ID, err)
	}

	Log.Info.Println(color.RedString("G:%v | C:%v > ACTION/DELETE > %v: %v", m.GuildID, m.ChannelID, msg.UserID, msg.Contents))
}
