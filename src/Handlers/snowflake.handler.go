package Handlers

import (
	"fmt"
	"math/big"
	"time"

	"antegr.al/chatanium-bot/v1/src/Log"
	"github.com/bwmarrin/discordgo"
)

func HandleSnowflake2time(s *discordgo.Session, i *discordgo.InteractionCreate) {
	id := i.ApplicationCommandData().Options[0].StringValue()

	Log.Verbose.Printf("Received Snowflake ID: %v", id)

	n := new(big.Int)
	n, ok := n.SetString(id, 10)
	if !ok {
		Log.Warn.Printf("Failed to convert string to bigint: %v", id)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Error: Failed to convert string to bigint. Please check your input.",
			},
		})
		return
	}

	timestamp := time.UnixMilli((n.Int64() >> 22) + 1420070400000)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("Timestamp: <t:%d> (Unix: `%d`)", timestamp.Unix(), timestamp.Unix()),
		},
	})
}
