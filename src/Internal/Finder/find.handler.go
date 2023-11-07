package Finder

import (
	"fmt"

	"antegr.al/chatanium-bot/v1/src/Database"
	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Log"

	"github.com/bwmarrin/discordgo"
)

func HandleFind(database *db.PrismaClient) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		id := i.ApplicationCommandData().Options[0].StringValue()
		isExpose := i.ApplicationCommandData().Options[1].BoolValue()

		flag := discordgo.MessageFlagsEphemeral

		if isExpose == true {
			flag = 0 // No Flag (omitempty)
		}

		Find := finder{
			session:  s,
			req:      i,
			database: database,
			flag:     flag,
		}

		switch id[0] {
		case 'C':
			// Channel
			Find.Channel()
			break
		case 'M':
			// Message
			Find.Message()
			break
		case 'G':
			// Guild
			break
		case 'U':
			// User
			break
		default:
			// Unknown
			Log.Warn.Printf("Invaild ID Prefix: %v", id[0])
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Error: Invaild ID Prefix. Please check your input.",
				},
			})
			break
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   flag,
				Content: "pong!",
			},
		})
	}
}

type finder struct {
	session  *discordgo.Session
	req      *discordgo.InteractionCreate
	database *db.PrismaClient
	flag     discordgo.MessageFlags
}

func (t *finder) Channel() {
}

func (t *finder) Message() {
	msg := Database.GetMessageInfo(t.req.GuildID, t.req.ChannelID, t.database)

	if msg == nil {
		t.session.InteractionRespond(t.req.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Error: Cannot find message. Please check your input.",
			},
		})
		return
	}

	content, err := msg.Contents()
	if err != true {
		t.session.InteractionRespond(t.req.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Error: Cannot retrieve message content.",
			},
		})
		return
	}

	t.session.InteractionRespond(t.req.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   t.flag,
			Content: fmt.Sprintf("**Message ID: %v**\nMessage Content: %v", msg.MessageID, content),
		},
	})
}
