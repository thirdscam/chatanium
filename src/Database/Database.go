package Database

import (
	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Util/Log"
)

// Database is a struct that contains the database client.
type Database struct {
	Client *db.PrismaClient
}

// Establish database connection. and, database must connected before start modules.
func (t *Database) Start() {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		Log.Error.Fatalf("Failed to connect to database: %v", err)
	}
	Log.Info.Println("Connected to database.")
	t.Client = client
}

// Close Database connection. must be called after all modules are shutdown.
func (t *Database) Shutdown() {
	Log.Verbose.Println("Shutting down database connection...")
	if err := t.Client.Disconnect(); err != nil {
		Log.Error.Panicf("Cannot close database connection: %v", err)
	}
	Log.Verbose.Println("Successfully closed!")
}
