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

	// if received a interrupt signal (CTRL+C), shutdown.
	signal.Notify(singal, os.Interrupt)
}
