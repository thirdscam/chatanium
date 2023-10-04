package main

import (
	"flag"
	"os"

	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Ignite"
	"antegr.al/chatanium-bot/v1/src/Log"
	"github.com/bwmarrin/discordgo"
)

var (
	Token       string
	LoggingMode int
)

func main() {
	flag.StringVar(&Token, "token", "", "Address to proxy")
	flag.IntVar(&LoggingMode, "logging-mode", 4, "Logging mode")

	// Parse the flags and init the logger
	flag.Parse()
	Log.Init(LoggingMode)

	Log.Info.Println("Chatanium: Scalable Bot Management System")
	Log.Info.Println("Press CTRL+C to shutdown.")

	// Create the session
	client := getClient("MTE1NDc4NTkzOTM5OTM3Njk2Ng.GEwjcR.Bc5uPjRJ1ceE8jtkqk3P4iLtCpbPIqx5Gq8brE")

	// Create a channel to receive OS signals
	interrupt := make(chan os.Signal)

	// Create a channel for disconnecting the database
	dbConn := make(chan *db.PrismaClient, 1)

	// Ignite Backend (Discord Bot, Status Page, etc.)
	go Ignite.Discord(interrupt, client, dbConn)

	// Wait for a signal to shutdown
	<-interrupt
	shutdown(interrupt, client, dbConn)
}

func shutdown(Singal chan os.Signal, Client *discordgo.Session, database chan *db.PrismaClient) {
	Log.Info.Println("Shutting down...")

	// Close the client
	if err := Client.Close(); err != nil {
		Log.Error.Panicf("Cannot close discord connection: %v", err)
	}
	Log.Verbose.Println("Discord connection closed.")

	// Close the database connection
	var dbConn *db.PrismaClient
	dbConn = <-database // Wait for the database to connect
	if err := dbConn.Prisma.Disconnect(); err != nil {
		Log.Error.Panicf("Cannot close database connection: %v", err)
	}
	Log.Verbose.Println("Database connection closed.")

	Log.Info.Println("Successfully shutdown.")
	os.Exit(0)
}

func getClient(token string) *discordgo.Session {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		Log.Error.Fatalln("Failed to create discord session: ", err)
	}

	return discord
}
