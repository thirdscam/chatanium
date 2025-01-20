package Database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	db "github.com/thirdscam/chatanium/src/Database/Internal"
	util "github.com/thirdscam/chatanium/src/Util"
	"github.com/thirdscam/chatanium/src/Util/Log"
)

func RegisterGuild(client *discordgo.Session, dbconn *sql.DB, queries *db.Queries, id, ownerID string) {
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
	Log.Verbose.Printf("G:%s > Found owner username: U:%s (%s)", id, ownerID, ownerUsername)

	// Get guild information
	g, err := client.Guild(id)
	if err != nil {
		Log.Error.Fatalf("G:%v > Failed to search guild: %v", id, err)
	}

	// Database: Begin transaction
	tx, err := dbconn.Begin()
	if err != nil {
		Log.Error.Fatalf("G:%s > Failed to begin transaction: %v", id, err)
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)

	// Database: insert user (Guild Owner)
	if err := qtx.InsertUser(context.Background(), db.InsertUserParams{
		ID:        util.Str2Int64(ownerID),
		Username:  ownerUsername,
		CreatedAt: st.JoinedAt,
	}); err != nil {
		Log.Error.Fatalf("G:%s > Failed to register guild (at INSERT_USER): %v", id, err)
		return
	}

	// Database: insert guild
	if err := qtx.InsertGuild(context.Background(), db.InsertGuildParams{
		ID:      util.Str2Int64(id),
		Name:    g.Name,
		OwnerID: util.Str2Int64(ownerID),
	}); err != nil {
		Log.Error.Fatalf("G:%s > Failed to register guild (at INSERT_GUILD): %v", id, err)
		return
	}

	// Database: insert member (Guild Owner)
	if err := qtx.InsertGuildUser(context.Background(), db.InsertGuildUserParams{
		UserID:    util.Str2Int64(ownerID),
		GuildID:   util.Str2Int64(id),
		CreatedAt: time.Now(),
	}); err != nil {
		Log.Error.Fatalf("G:%s > Failed to insert user (at INSERT_MEMBER): %v", id, err)
		return
	}

	// Insert each channel in guild
	var cntNewChannel int
	Log.Verbose.Printf("G:%s > Inserting channels... (%d channels)", g.ID, len(g.Channels))
	for _, v := range g.Channels {
		// Database: insert channel
		if err := qtx.InsertChannel(context.Background(), db.InsertChannelParams{
			ID:      util.Str2Int64(v.ID),
			GuildID: util.Str2Int64(g.ID),
			Name:    v.Name,
		}); err != nil {
			Log.Warn.Printf("G:%s | C:%s > Failed to insert channel: %v", g.ID, v.ID, err)
		}

		cntNewChannel++
	}

	if err := tx.Commit(); err != nil {
		Log.Warn.Printf("G:%s > Failed to commit transaction: %v", g.ID, err)
		return
	}

	// If new channels inserted
	if cntNewChannel > 0 {
		Log.Verbose.Printf("G:%s > Inserted %d new channels.", g.ID, cntNewChannel)
	}

	// If all channels already inserted
	Log.Verbose.Printf("G:%s > Done.", g.ID)
}

// InsertUser inserts user information into the database
func InsertUser(queries *db.Queries, uid, username string) {
	// Check if user already exists
	_, err := queries.GetUser(context.Background(), util.Str2Int64(uid))
	if err == nil {
		// User already exists
		return
	} else if err != sql.ErrNoRows {
		Log.Error.Fatalf("Failed to find user: %v", err)
	}

	// Insert user
	err = queries.InsertUser(context.Background(), db.InsertUserParams{
		ID:        util.Str2Int64(uid),
		Username:  username,
		CreatedAt: time.Now(),
	})
	if err != nil {
		Log.Error.Fatalf("Failed to insert user: %v", err)
	}
	Log.Verbose.Printf("U:%s > User insert completed.", uid)
}

// InsertMember inserts member information into the database
func InsertMember(queries *db.Queries, uid, gid, nickname string) bool {
	// Check if guild user already exists
	_, err := queries.GetGuildUser(context.Background(), db.GetGuildUserParams{
		UserID:  util.Str2Int64(uid),
		GuildID: util.Str2Int64(gid),
	})
	if err == nil {
		// Guild user already exists
		return false
	} else if err != sql.ErrNoRows {
		log.Fatalf("Failed to find guild user: %v", err)
	}

	log.Printf("G:%s | U:%s > Adding member... (%s)", gid, uid, nickname)

	// Insert guild user
	newUUID := uuid.New().String()
	err = queries.InsertGuildUser(context.Background(), db.InsertGuildUserParams{
		Uuid:      newUUID,
		GuildID:   util.Str2Int64(gid),
		UserID:    util.Str2Int64(uid),
		CreatedAt: time.Now(),
		Nickname:  nickname,
	})
	if err != nil {
		log.Fatalf("Failed to insert guild user: %v", err)
	}
	log.Printf("G:%s | U:%s > Member insert completed. (%s)", gid, uid, nickname)
	return true
}

// // InsertGuild inserts guild information into the database
// func InsertGuild(dbConn *sql.DB, gid, name, ownerUid string) {
// 	ctx := context.Background()
// 	queries := db.New(dbConn)

// 	// Check if guild already exists
// 	_, err := queries.GetGuild(ctx, util.Str2Int64(gid))
// 	if err == nil {
// 		// Guild already exists
// 		return
// 	} else if err != sql.ErrNoRows {
// 		log.Fatalf("Failed to find guild: %v", err)
// 	}

// 	// Insert guild
// 	err = queries.InsertGuild(ctx, db.InsertGuildParams{
// 		ID:      util.Str2Int64(gid),
// 		Name:    name,
// 		OwnerID: util.Str2Int64(ownerUid),
// 	})
// 	if err != nil {
// 		Log.Error.Fatalf("Failed to insert guild: %v", err)
// 	}
// 	Log.Verbose.Printf("G:%s > Guild insert completed.", gid)
// }

// // InsertChannel inserts channel information into the database
// func InsertChannel(dbConn *sql.DB, cid, gid, name, description string) bool {
// 	ctx := context.Background()
// 	queries := db.New(dbConn)

// 	// Check if channel already exists
// 	_, err := queries.GetChannel(ctx, util.Str2Int64(cid))
// 	if err == nil {
// 		// Channel already exists
// 		return false
// 	} else if err != sql.ErrNoRows {
// 		Log.Error.Fatalf("Failed to find channel: %v", err)
// 	}

// 	// Insert channel
// 	err = queries.InsertChannel(ctx, db.InsertChannelParams{
// 		ID:        util.Str2Int64(cid),
// 		GuildID:   util.Str2Int64(gid),
// 		Name:      name,
// 		CreatedAt: time.Now(),
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed to insert channel: %v", err)
// 	}
// 	log.Printf("G:%s | C:%s > Channel insert completed.", gid, cid)
// 	return true
// }

// // InsertMember inserts member information into the database
// func InsertMember(dbConn *sql.DB, uid, gid, nickname string) bool {
// 	ctx := context.Background()
// 	queries := db.New(dbConn)

// 	// Check if guild user already exists
// 	_, err := queries.GetGuildUser(ctx, db.GetGuildUserParams{
// 		UserID:  util.Str2Int64(uid),
// 		GuildID: util.Str2Int64(gid),
// 	})
// 	if err == nil {
// 		// Guild user already exists
// 		return false
// 	} else if err != sql.ErrNoRows {
// 		log.Fatalf("Failed to find guild user: %v", err)
// 	}

// 	log.Printf("G:%s | U:%s > Adding member... (%s)", gid, uid, nickname)

// 	// Insert guild user
// 	newUUID := uuid.New().String()
// 	err = queries.InsertGuildUser(ctx, db.InsertGuildUserParams{
// 		UUID:      newUUID,
// 		GuildID:   util.Str2Int64(gid),
// 		UserID:    util.Str2Int64(uid),
// 		CreatedAt: time.Now(),
// 		Nickname:  nickname,
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed to insert guild user: %v", err)
// 	}
// 	log.Printf("G:%s | U:%s > Member insert completed. (%s)", gid, uid, nickname)
// 	return true
// }

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
