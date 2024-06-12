package Database

import (
	"context"
	"errors"
	"time"

	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	util "antegr.al/chatanium-bot/v1/src/Util"
	"antegr.al/chatanium-bot/v1/src/Util/Log"
)

func InsertUser(database *db.PrismaClient, uid string, username string) {
	ctx := context.Background()
	// Log.Verbose.Printf("U:%s (%s) > Adding user...", uid, username)

	Users := db.Users

	_, err := database.Users.FindUnique(
		Users.ID.Equals(util.StringToBigint(uid)),
	).Exec(ctx)
	if err == nil {
		// Log.Verbose.Printf("U:%s (%s) > User already exists.", uid, username)
		return
	} else if !errors.Is(err, db.ErrNotFound) {
		Log.Error.Fatalf("U:%s (%s) > Failed to find user: %v", uid, username, err)
	}

	_, err = database.Users.CreateOne(
		Users.ID.Set(util.StringToBigint(uid)),
		Users.Username.Set(username),
		Users.CreatedAt.Set(time.Now()),
	).Exec(ctx)
	if err != nil {
		Log.Error.Panicf("U:%s (%s) > Failed to register user: %v", uid, username, err)
		// Log.Verbose.Printf("U:%s (%s) > User already exists.", uid, username)
		return
	}
}
