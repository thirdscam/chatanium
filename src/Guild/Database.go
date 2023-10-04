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

func RegisterDatabase(database *db.PrismaClient, g *discordgo.GuildCreate) {
	Log.Verbose.Printf("%s > Adding database...", g.ID)

	Database.UpsertUser(database, g.OwnerID, g.OwnerID)
	UpsertGuild(database, g.ID, g.Name, g.OwnerID)

	for _, v := range g.Channels {
		UpsertChannel(database, v.ID, g.ID, v.Name, v.Topic)
	}
}

func UpsertGuild(database *db.PrismaClient, gid, name, ownerUid string) {
	ctx := context.Background()

	// Database Task: Upsert guild
	Guild := db.Guilds
	_, err := database.Guilds.UpsertOne(Guild.ID.Equals(util.StringToBigint(gid))).Create(
		Guild.ID.Set(util.StringToBigint(gid)),
		Guild.Name.Set(name),
		Guild.Users.Link(db.Users.ID.Equals(util.StringToBigint(ownerUid))),
	).Exec(ctx)
	if err != nil {
		Log.Error.Fatalf("Failed to upsert guild: %v", err)
	}
}

func UpsertChannel(database *db.PrismaClient, cid, gid, name, description string) {
	ctx := context.Background()

	// Database Task: Upsert channel
	Channel := db.Channels
	_, err := database.Channels.UpsertOne(Channel.ID.Equals(util.StringToBigint(cid))).Create(
		Channel.ID.Set(util.StringToBigint(cid)),
		Channel.Name.Set(name),
		Channel.CreatedAt.Set(time.Now()),
		Channel.Guilds.Link(db.Guilds.ID.Equals(util.StringToBigint(gid))),
	).Exec(ctx)
	if err != nil {
		Log.Error.Fatalf("Failed to upsert channel: %v", err)
	}
}

func UpsertUser(database *db.PrismaClient, uid string, username string) {
	ctx := context.Background()

	// Database Task: Upsert user (Guild Member)
	Users := db.Guildusers

	user, err := database.Guildusers.FindFirst(Users.UserID.Equals(util.StringToBigint(uid))).Exec(ctx)
	if errors.Is(err, db.ErrNotFound) {
		_, err = database.Guildusers.CreateOne(
			Users.UUID.Set(""),
			Users.Username.Set(username),
			Users.CreatedAt.Set(time.Now()),
		).Exec(ctx)
		if err != nil {
			Log.Error.Fatalf("Failed to upsert user: %v", err)
		}
	} else if err != nil {
		Log.Error.Fatalf("Failed to find user: %v", err)
	}

	_, err = database.Guildusers.UpsertOne(Users.UserID.Equals(util.StringToBigint(uid))).Create(
		Users.ID.Set(util.StringToBigint(uid)),
		Users.Username.Set(username),
		Users.CreatedAt.Set(time.Now()),
	).Exec(ctx)
	if err != nil {
		Log.Error.Fatalf("Failed to upsert user: %v", err)
	}
}
