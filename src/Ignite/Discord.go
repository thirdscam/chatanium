package Ignite

import (
	"os"
	"os/signal"

	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Guild"
	"antegr.al/chatanium-bot/v1/src/Log"
	"github.com/bwmarrin/discordgo"
)

func Discord(singal chan os.Signal, client *discordgo.Session, db *db.PrismaClient) {
	Log.Info.Println("Starting Bot...")

	// Open the connection from discord
	if err := client.Open(); err != nil {
		Log.Error.Fatalf("Cannot open connection: %v", err)
	}

	// Getting Token infomation
	client.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		Log.Verbose.Printf("WS/READY > Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	client.AddHandler(func(s *discordgo.Session, r *discordgo.GuildDelete) {
		Log.Verbose.Printf("WS/READY > Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Handle all guilds
	Guild.Handle(client, db)

	// if received a interrupt signal (CTRL+C), shutdown.
	signal.Notify(singal, os.Interrupt)
}
