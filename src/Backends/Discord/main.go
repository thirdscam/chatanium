package discord

import (
	"strings"

	"antegr.al/chatanium-bot/v1/src/Backends/Discord/Guild"
	"antegr.al/chatanium-bot/v1/src/Database"
	module "antegr.al/chatanium-bot/v1/src/Module"
	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/bwmarrin/discordgo"
)

type Backend struct {
	Token    string
	Database *Database.DB
	client   *discordgo.Session
}

// Start Discord backend.
// This function is Entry Point of Discord backend.
func (t *Backend) Start() {
	// Remove "Bot " prefix from token
	t.Token = strings.TrimPrefix(t.Token, "Bot ")

	client, err := discordgo.New("Bot " + t.Token)
	if err != nil {
		Log.Error.Fatalf("Failed to create Discord session: %v", err)
	}

	if err := client.Open(); err != nil {
		Log.Error.Fatalf("Failed to close Discord session: %v", err)
	}

	// Get modules
	moduleManager := module.ModuleManager{
		Identifier: "discord",
	}
	moduleManager.Start()

	t.client = client

	Guild.Handle(client, t.Database)
	// DM.Handle(client, db) // TODO: DM (Direct Message) Handler
}

// Shutdown Discord backend.
func (t *Backend) Shutdown() {
	Log.Verbose.Println("Shutting down Discord backend...")
	err := t.client.Close()
	if err != nil {
		Log.Error.Fatalf("Failed to close Discord session: %v", err)
	}
	Log.Verbose.Println("Successfully closed!")
}
