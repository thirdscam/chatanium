package discord

import (
	"strings"

	"antegr.al/chatanium-bot/v1/src/Backends/Discord/Guild"
	"antegr.al/chatanium-bot/v1/src/Backends/Discord/Interface/Slash"
	"antegr.al/chatanium-bot/v1/src/Backends/Discord/Module"
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

	// 1. Create Discord session
	client, err := discordgo.New("Bot " + t.Token)
	if err != nil {
		Log.Error.Fatalf("Failed to create Discord session: %v", err)
	}

	// 2. Get modules
	moduleMgr := struct{ module.ModuleManager }{
		ModuleManager: module.ModuleManager{
			Identifier: "discord",
		},
	}
	SlashGuildMgr := Slash.Guild{
		Client: client,
	}

	// 3. Load modules
	moduleMgr.Load()
	CommandSet := Slash.Commands{}

	// TODO(Feature): Improved command mapping (ACL, etc.)
	for _, v := range moduleMgr.Modules {
		module := Module.ConvertDiscordModule(v)
		for k, v := range module.Commands {
			CommandSet[k] = v
		}
	}

	Slash.CommandMap = CommandSet

	// 4. Start modules and Interface managers
	moduleMgr.Start()
	SlashGuildMgr.Start()
	t.client = client

	// 4. Add handlers
	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		SlashGuildMgr.OnGuildCreated(g.ID) // Register slash commands
	})

	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildDelete) {
		SlashGuildMgr.OnGuildDeleted(g.ID) // Remove slash commands
	})
	Guild.Handle(client, t.Database)
	// DM.Handle(client, db) // TODO: DM (Direct Message) Handler

	// 5. Start Discord session
	if err := client.Open(); err != nil {
		Log.Error.Fatalf("Failed to start Discord session: %v", err)
	}
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
