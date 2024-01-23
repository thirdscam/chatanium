package discord

import (
	"antegr.al/chatanium-bot/v1/src/Backends/Discord/Guild"
	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"github.com/bwmarrin/discordgo"
)

// Start Discord backend.
// This function is Entry Point of Discord backend.
func Start(client *discordgo.Session, db *db.PrismaClient) {
	Guild.Handle(client, db)
	// DM.Handle(client, db) // TODO: DM (Direct Message) Handler
}
