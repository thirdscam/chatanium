package Ignite

import (
	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Guild"
	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Database *db.PrismaClient
	Token    string
	client   *discordgo.Session
}

// Start the discord bot backend
func (t *Discord) Start() {
	Log.Info.Println("Starting Bot...")

	client, err := discordgo.New("Bot " + t.Token)
	if err != nil {
		Log.Error.Fatalln("Failed to create discord session: ", err)
	}

	t.client = client

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
	go Guild.Handle(client, t.Database)
}

// Close the connection from discord
func (t *Discord) Shutdown() {
	Log.Verbose.Println("Shutting down discord connection...")
	if err := t.client.Close(); err != nil {
		Log.Error.Panicf("Cannot close discord connection: %v", err)
	}
	Log.Verbose.Println("Closed successfully!")
}
