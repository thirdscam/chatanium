package Ignite

import (
	"os"

	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Log"
	"github.com/bwmarrin/discordgo"
)

func Shutdown(Singal chan os.Signal, Client *discordgo.Session, database *db.PrismaClient) {
	<-Singal // Wait for a signal
	Log.Info.Println("Shutting down...")

	// Close the client
	if err := Client.Close(); err != nil {
		Log.Error.Panicf("Cannot close discord connection: %v", err)
	}
	Log.Verbose.Println("Discord connection closed.")

	// Close the database connection
	if err := database.Prisma.Disconnect(); err != nil {
		Log.Error.Panicf("Cannot close database connection: %v", err)
	}
	Log.Verbose.Println("Database connection closed.")

	Log.Info.Println("Successfully shutdown.")
	os.Exit(0)
}
