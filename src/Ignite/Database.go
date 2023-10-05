package Ignite

import (
	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Log"
)

func DB() *db.PrismaClient {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		Log.Error.Fatalf("Failed to connect to database: %v", err)
	}

	return client
}
