package Database

import (
	"context"
	"time"

	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Log"
	"antegr.al/chatanium-bot/v1/src/util"
)

func Get() (error, *db.PrismaClient) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return err, nil
	}

	return nil, client
}

func UpsertUser(database *db.PrismaClient, uid string, username string) {
	ctx := context.Background()

	// Database Task: Upsert user (Register User)
	Users := db.Users
	_, err := database.Users.UpsertOne(Users.ID.Equals(util.StringToBigint(uid))).Create(
		Users.ID.Set(util.StringToBigint(uid)),
		Users.Username.Set(username),
		Users.CreatedAt.Set(time.Now()),
	).Exec(ctx)
	if err != nil {
		Log.Error.Fatalf("Failed to upsert user: %v", err)
	}
}
