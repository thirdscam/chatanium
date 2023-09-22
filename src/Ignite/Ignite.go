package Ignite

import (
	"os"
	"os/signal"

	"antegr.al/chatanium-bot/v1/src/Guild"
	"antegr.al/chatanium-bot/v1/src/Handlers"
	"antegr.al/chatanium-bot/v1/src/Log"
	"antegr.al/chatanium-bot/v1/src/Schema"
	"github.com/bwmarrin/discordgo"
)

func Discord(singal chan os.Signal, client *discordgo.Session) {
	Log.Info.Println("Starting Bot...")

	// Open the connection from discord
	if err := client.Open(); err != nil {
		Log.Error.Fatalf("Cannot open connection: %v", err)
	}

	// Getting Token infomation
	client.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		Log.Info.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// TODO: Database (Save Guild ID, etc.)
	var GuildCmds []Guild.Commands

	// Register all commands from all guilds
	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		Log.Verbose.Printf("Joined Guild: %v (%v)", g.Name, g.ID)

		AllowedModules := []string{"ping"}

		Guild := Guild.Commands{
			Schema:   Schema.GetAllowedOnly(AllowedModules),
			Handlers: Handlers.GetAllowedOnly(AllowedModules),
			Client:   client,
			GuildID:  g.ID,
		}

		// Register commands from guild
		Guild.RegisterHandlers()
		Guild.RegisterSchema()

		GuildCmds = append(GuildCmds, Guild)
	})

	// Handle all messages from all guilds
	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == client.State.User.ID {
			return
		}

		Log.Verbose.Printf("(%v -> %v) %v: %v", m.GuildID, m.ChannelID, m.Author, m.Content)
	})

	// if received a interrupt signal (CTRL+C), shutdown.
	signal.Notify(singal, os.Interrupt)
}
