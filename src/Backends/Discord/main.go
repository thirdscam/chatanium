package discord

import (
	"antegr.al/chatanium-bot/v1/src/Backends/Discord/Guild"
	"antegr.al/chatanium-bot/v1/src/Database"
	"github.com/bwmarrin/discordgo"
)

// Start Discord backend.
// This function is Entry Point of Discord backend.
func Start(client *discordgo.Session, db *Database.DB) {
	Guild.Handle(client, db)
	// DM.Handle(client, db) // TODO: DM (Direct Message) Handler
}
