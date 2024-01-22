package Database

import (
	"context"
	"errors"
	"time"

	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	util "antegr.al/chatanium-bot/v1/src/Util"
	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func RegisterGuild(client *discordgo.Session, database *db.PrismaClient, id, ownerID string) {
	Log.Verbose.Printf("G:%s > Adding database...", id)

	// Search owner username by owner id
	st, err := client.GuildMember(id, ownerID)
	if err != nil {
		Log.Error.Fatalf("G:%v | U:%v > Failed to get member: %v", id, ownerID, err)
	}
	if st.User.Username == "" {
		Log.Error.Fatalf("G:%s > Failed to find owner username from id", id)
	}

	// Username of the guild owner
	ownerUsername := st.User.Username
	Log.Verbose.Printf("G:%s > Found owner username: %s (%s)", id, ownerUsername, ownerID)

	// Database: insert user (Guild Owner)
	InsertUser(database, ownerID, ownerUsername)

	g, err := client.Guild(id)
	if err != nil {
		Log.Error.Fatalf("G:%v > Failed to search guild: %v", id, err)
	}

	// Database: insert guild
	InsertGuild(database, id, g.Name, ownerID)

	// Database: insert member (Guild Owner)
	InsertMember(database, ownerID, id, ownerUsername)

	// insert each channels in guild
	var cntNewChannel int
	for _, v := range g.Channels {
		// Database: insert channel
		isInserted := InsertChannel(database, v.ID, g.ID, v.Name, v.Topic)
		if isInserted {
			cntNewChannel++
		}
	}

	// if new channels inserted
	if cntNewChannel > 0 {
		Log.Verbose.Printf("G:%s > Inserted %d new channels.", g.ID, cntNewChannel)
		return
	}

	// if all channels already inserted
	Log.Verbose.Printf("G:%s > All channels already exists.", g.ID)
}

// Insert guild information to database
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

// Insert channel information to database
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

// Insert user information to database.
// (Member not equals to user. Member is a user who joined the guild)
func InsertMember(database *db.PrismaClient, uid, gid, nickname string) bool {
	ctx := context.Background()
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

	Log.Verbose.Printf("G:%s | U:%s > Adding member... (%s)", gid, uid, nickname)

	// Database Task: Insert user (Guild Member)
	newUuid := uuid.New().String() // make new uuid
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

// TODO(feature): Updata guild information to database
func UpdateGuild() {
}

// TODO(feature): Remove guild information from database
func RemoveGuild() {
}

// TODO(feature): Update Member information from database
// (Member not equals to user. Member is a user who joined the guild)
func UpdateMember() {
}
