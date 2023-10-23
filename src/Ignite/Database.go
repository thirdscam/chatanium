package Ignite

import (
	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Log"
)

type DB struct {
	client *db.PrismaClient
}

func (t *DB) Start() *db.PrismaClient {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		Log.Error.Fatalf("Failed to connect to database: %v", err)
	}

	t.client = client

	return client
}

func (t *DB) Shutdown() {
	Log.Verbose.Println("Shutting down database connection...")
	if err := t.client.Disconnect(); err != nil {
		Log.Error.Panicf("Cannot close database connection: %v", err)
	}
	Log.Verbose.Println("Closed successfully!")
}
