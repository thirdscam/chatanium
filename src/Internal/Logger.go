package Internal

import (
	"context"
	"errors"
	"fmt"

	"antegr.al/chatanium-bot/v1/src/Database"
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

		// TODO: global_name (default nickname) unsupported issue
		//
		// Member.Nick에서 Profile by Guild에 대한 Nickname을 가져올 수 있음.
		// 이는 Profile by Guild이므로, Default Value는 ""임.
		//
		// 따라서, default nickname을 User.global_name에서 가져와야 함.
		// 단, 이 기능은 현재 discordgo에서는 지원하지 않음.
		//
		// v0.27.1 (Latest, 2023/03): https://pkg.go.dev/github.com/bwmarrin/discordgo#User
		// Development (Unreleased, 2023/07): https://github.com/bwmarrin/discordgo/blame/master/user.go#L53

		// TODO: Nickname failback
		// 현재 모든 유저에 대한 Nickname 수집을 보장 할 수 없음.
		// 따라서, Nickname을 null로 설정하고, failback 구현.

		// Search owner nickname
		u, _ := client.GuildMember(m.GuildID, m.Author.ID)
		_ = u

		// st.ni
		st, err := client.GuildMember(m.GuildID, m.Author.ID)
		// st, err := s.User(m.Author.ID)
		if err != nil {
			Log.Error.Fatalf("G:%v | U:%v > Failed to get member: %v", m.GuildID, m.Author.ID, err)
		}
		fmt.Printf("%v", st)
		// if st.Flags == "" {
		// 	Log.Error.Fatalf("G:%s > Failed to find owner username from id", m.GuildID)
		// }

		createMessage(m, database)
		Log.Verbose.Printf("G:%v | C:%v > %v: %v", m.GuildID, m.ChannelID, m.Author.Username, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == client.State.User.ID {
			return
		}

		Log.Verbose.Printf("G:%v | C:%v > Update M:%v > %v", m.GuildID, m.ChannelID, m.Message.ID, m.Content)
	})

	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageDelete) {
		Log.Verbose.Printf("G:%v | C:%v > Delete M:%v", m.GuildID, m.ChannelID, m.ID)

		msg := getMessageInfo(m.GuildID, m.ID, database)
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

		Log.Info.Printf("G:%v | C:%v > Delete Message > %v: %v", m.GuildID, m.ChannelID, msg.Users().ID, content)
	})
}

func MemberLogger(client *discordgo.Session, dbClient *db.PrismaClient) {
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

func createMessage(m *discordgo.MessageCreate, database *db.PrismaClient) {
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

func getMessageInfo(gid, mid string, database *db.PrismaClient) *db.MessagesModel {
	msg, err := database.Messages.FindUnique(
		db.Messages.MessageID.Equals(util.StringToBigint(mid)),
	).Exec(
		context.Background(),
	)
	if err == nil {
		return nil
	} else if !errors.Is(err, db.ErrNotFound) {
		Log.Error.Printf("G:%s | M:%s > Failed to find message : %v", gid, mid, err)
	}

	return msg
}
