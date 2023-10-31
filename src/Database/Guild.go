package Database

import (
	"context"
	"errors"
	"time"

	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Log"
	util "antegr.al/chatanium-bot/v1/src/Util"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func RegisterGuild(client *discordgo.Session, database *db.PrismaClient, g *discordgo.GuildCreate) {
	Log.Verbose.Printf("G:%s > Adding database...", g.ID)

	// Search owner username
	st, err := client.GuildMember(g.ID, g.OwnerID)
	if err != nil {
		Log.Error.Fatalf("G:%v | U:%v > Failed to get member: %v", g.ID, g.OwnerID, err)
	}
	if st.User.Username == "" {
		Log.Error.Fatalf("G:%s > Failed to find owner username from id", g.ID)
	}
	OwnerUsername := st.User.Username
	Log.Verbose.Printf("G:%s > Found owner username: %s (%s)", g.ID, OwnerUsername, g.OwnerID)

	// Database insert user (Guild Owner)
	InsertUser(database, g.OwnerID, OwnerUsername)

	// Database insert guild
	InsertGuild(database, g.ID, g.Name, g.OwnerID)

	// Database insert member (Guild Owner)
	InsertMember(database, g.OwnerID, g.ID, OwnerUsername)

	var Inserted int
	for _, v := range g.Channels {
		// Database insert channel
		isInserted := InsertChannel(database, v.ID, g.ID, v.Name, v.Topic)
		if isInserted {
			Inserted++
		}
	}

	if Inserted > 0 {
		Log.Verbose.Printf("G:%s > Inserted %d new channels.", g.ID, Inserted)
		return
	}

	Log.Verbose.Printf("G:%s > All channels already exists.", g.ID)
}

func InsertGuild(database *db.PrismaClient, gid, name, ownerUid string) {
	ctx := context.Background()

	// Database Task: Upsert guild
	Guild := db.Guilds
	_, err := database.Guilds.FindUnique(
		Guild.ID.Equals(util.StringToBigint(gid)),
	).Exec(ctx)
	if err == nil {
		// Log.Verbose.Printf("G:%s > Guild already exists.", gid)
		return
	} else if !errors.Is(err, db.ErrNotFound) {
		Log.Error.Fatalf("G:%s > Failed to find guild: %v", gid, err)
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

func InsertChannel(database *db.PrismaClient, cid, gid, name, description string) bool {
	ctx := context.Background()

	// Database Task: Insert channel
	Channel := db.Channels
	_, err := database.Channels.FindUnique(
		Channel.ID.Equals(util.StringToBigint(cid)),
	).Exec(ctx)
	if err == nil {
		// Log.Verbose.Printf("G:%s | C:%s > Channel already exists.", gid, cid)
		return false
	} else if !errors.Is(err, db.ErrNotFound) {
		Log.Error.Fatalf("G:%s | C:%s > Failed to find guild: %v", gid, cid, err)
	}

	_, err = database.Channels.CreateOne(
		Channel.ID.Set(util.StringToBigint(cid)),
		Channel.Name.Set(name),
		Channel.CreatedAt.Set(time.Now()),
		Channel.Guilds.Link(db.Guilds.ID.Equals(util.StringToBigint(gid))),
	).Exec(ctx)
	if err != nil {
		Log.Error.Fatalf("Failed to insert channel: %v", err)
	}

	Log.Verbose.Printf("G:%s | C:%s > Channel insert completed.", gid, cid)
	return true
}

func InsertMember(database *db.PrismaClient, uid, gid, nickname string) bool {
	ctx := context.Background()
	Log.Verbose.Printf("G:%s | U:%s > Adding member... (%s)", gid, uid, nickname)

	// Database task: Check exists member in guild
	Guilduser := db.Guildusers
	_, err := database.Guildusers.FindFirst(
		Guilduser.UserID.Equals(util.StringToBigint(uid)),
		Guilduser.GuildID.Equals(util.StringToBigint(gid)),
	).Exec(ctx)
	if err == nil {
		// Log.Verbose.Printf("G:%s > Guild user already exists.", uid)
		return false
	} else if !errors.Is(err, db.ErrNotFound) {
		Log.Error.Fatalf("G:%s > Failed to find guild user: %v", uid, err)
	}

	// Database Task: Insert user (Guild Member)
	newUuid := uuid.New().String()
	a, err := database.Guildusers.CreateOne(
		Guilduser.CreatedAt.Set(time.Now()),
		Guilduser.Nickname.Set(nickname),
		Guilduser.UUID.Set(newUuid),
		Guilduser.Guilds.Link(db.Guilds.ID.Equals(util.StringToBigint(gid))),
		Guilduser.Users.Link(db.Users.ID.Equals(util.StringToBigint(uid))),
	).Exec(ctx)
	if err != nil {
		Log.Error.Fatalf("Failed to insert guild user: %v", err)
	}

	Log.Verbose.Printf("G:%v | U:%v > Registered! (%s)", a.GuildID, a.UserID, a.Nickname)

	return true
}
