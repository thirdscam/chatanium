package Guild

import (
	"context"
	"errors"
	"time"

	"antegr.al/chatanium-bot/v1/src/Database"
	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Log"
	"antegr.al/chatanium-bot/v1/src/util"
	"github.com/bwmarrin/discordgo"
)

func RegisterDatabase(client *discordgo.Session, database *db.PrismaClient, g *discordgo.GuildCreate) {
	Log.Verbose.Printf("%s > Adding database...", g.ID)

	OwnerUsername := SearchUsernameByUID(client, g.OwnerID, g.ID)
	if OwnerUsername == "" {
		Log.Error.Fatalf("%s > Failed to find owner username from id", g.ID)
	}

	Log.Verbose.Printf("%s > Found owner username: %s (%s)", g.ID, OwnerUsername, g.OwnerID)

	Database.InsertUser(database, g.OwnerID, OwnerUsername)
	InsertGuild(database, g.ID, g.Name, g.OwnerID)

	for _, v := range g.Channels {
		UpsertChannel(database, v.ID, g.ID, v.Name, v.Topic)
	}
}

func InsertGuild(database *db.PrismaClient, gid, name, ownerUid string) {
	ctx := context.Background()

	// Database Task: Upsert guild
	Guild := db.Guilds
	_, err := database.Guilds.FindUnique(
		Guild.ID.Equals(util.StringToBigint(gid)),
	).Exec(ctx)
	if err == nil {
		Log.Verbose.Printf("%s > Guild already exists.", gid)
		return
	} else if !errors.Is(err, db.ErrNotFound) {
		Log.Error.Fatalf("%s > Failed to find guild: %v", gid, err)
	}

	_, err = database.Guilds.CreateOne(
		Guild.ID.Set(util.StringToBigint(gid)),
		Guild.Name.Set(name),
		Guild.Users.Link(db.Users.ID.Equals(util.StringToBigint(ownerUid))),
	).Exec(ctx)
	if err != nil {
		Log.Error.Fatalf("Failed to insert guild: %v", err)
	}

	Log.Verbose.Printf("%s > Guild insert completed.", gid)
}

func UpsertChannel(database *db.PrismaClient, cid, gid, name, description string) {
	ctx := context.Background()

	// Database Task: Upsert channel
	Channel := db.Channels
	_, err := database.Channels.FindUnique(
		Channel.ID.Equals(util.StringToBigint(cid)),
	).Exec(ctx)
	if err == nil {
		Log.Verbose.Printf("%s::%s > Channel already exists.", gid, cid)
		return
	} else if !errors.Is(err, db.ErrNotFound) {
		Log.Error.Fatalf("%s::%s > Failed to find guild: %v", gid, cid, err)
	}

	_, err = database.Channels.CreateOne(
		Channel.ID.Set(util.StringToBigint(cid)),
		Channel.Name.Set(name),
		Channel.CreatedAt.Set(time.Now()),
		Channel.Guilds.Link(db.Guilds.ID.Equals(util.StringToBigint(gid))),
	).Exec(ctx)
	if err != nil {
		Log.Error.Fatalf("Failed to upsert channel: %v", err)
	}

	Log.Verbose.Printf("%s::%s > Channel insert completed.", gid, cid)
}

// func UpsertUser(database *db.PrismaClient, uid, gid, username string) {
// 	ctx := context.Background()

// 	// Database Task: Upsert user (Guild Member)
// 	Users := db.Guildusers

// 	_, err := database.Guildusers.FindFirst(Users.UserID.Equals(util.StringToBigint(uid))).Exec(ctx)
// 	if errors.Is(err, db.ErrNotFound) {
// 		_, err = database.Guildusers.CreateOne(
// 			Users.UUID.Set(uuid.New().String()),
// 			Users.UserID.Set(util.StringToBigint(uid)),
// 			Users.CreatedAt.Set(time.Now()),
// 			,
// 		).Exec(ctx)
// 		if err != nil {
// 			Log.Error.Fatalf("Failed to upsert user: %v", err)
// 		}
// 	} else if err != nil {
// 		Log.Error.Fatalf("Failed to find user: %v", err)
// 	}

// 	_, err = database.Guildusers.UpsertOne(Users.UserID.Equals(util.StringToBigint(uid))).Create(
// 		Users.ID.Set(util.StringToBigint(uid)),
// 		Users.Username.Set(username),
// 		Users.CreatedAt.Set(time.Now()),
// 	).Exec(ctx)
// 	if err != nil {
// 		Log.Error.Fatalf("Failed to upsert user: %v", err)
// 	}
// }
