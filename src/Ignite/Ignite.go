package Ignite

import (
	"os"
	"os/signal"

	"antegr.al/chatanium-bot/v1/src/Guild"
	"antegr.al/chatanium-bot/v1/src/Log"
	"github.com/bwmarrin/discordgo"
)

func Discord(singal chan os.Signal, client *discordgo.Session, RegisteredGuildCmds *[]Guild.Commands) {
	Log.Info.Println("Starting Bot...")

	// Getting Token infomation
	client.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		Log.Info.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// TODO: Database (Save Guild ID, etc.)
	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		Log.Info.Printf("Joined Guild: %v (%v)", g.Name, g.ID)
		RegisteredGuildCmds = append(*RegisteredGuildCmds, Guild.Commands{
			Schema:   Guild.GetSchema(),
			Handlers: Guild.GetHandlers(),
			Client:   client,
			GuildID:  g.ID,
		})
	})

	// Register all commands from all guilds

	// if received a interrupt signal (CTRL+C), shutdown.
	signal.Notify(singal, os.Interrupt)
}
