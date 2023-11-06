package Handlers

// import (
// 	"antegr.al/chatanium-bot/v1/src/Database"
// 	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
// 	"antegr.al/chatanium-bot/v1/src/Log"
// 	"github.com/bwmarrin/discordgo"
// )

// func HandleFind(s *discordgo.Session, i *discordgo.InteractionCreate, Database *db.PrismaClient) {
// 	id := i.ApplicationCommandData().Options[0].StringValue()
// 	isExpose := i.ApplicationCommandData().Options[1].BoolValue()

// 	switch id[0] {
// 	case 'C':
// 		// Channel
// 		break
// 	case 'M':
// 		// Message
// 		break
// 	case 'G':
// 		// Guild
// 		break
// 	case 'U':
// 		// User
// 		break
// 	default:
// 		// Unknown
// 		Log.Warn.Printf("Invaild ID Prefix: %v", id[0])
// 		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
// 			Data: &discordgo.InteractionResponseData{
// 				Flags:   discordgo.MessageFlagsEphemeral,
// 				Content: "Error: Invaild ID Prefix. Please check your input.",
// 			},
// 		})
// 		break
// 	}

// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
// 		Data: &discordgo.InteractionResponseData{
// 			Content: "pong!",
// 		},
// 	})
// }

// func searchChannel(s *discordgo.Session, i *discordgo.InteractionCreate, database *db.PrismaClient) {
// }

// func searchMessage(s *discordgo.Session, i *discordgo.InteractionCreate, database *db.PrismaClient) {
// 	msg := Database.GetMessageInfo(i.GuildID, i.ChannelID, database)

// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
// 		Data: &discordgo.InteractionResponseData{
// 			Flags:   discordgo.MessageFlagsEphemeral,
// 			Content: "Error: Invaild ID Prefix. Please check your input.",
// 		},
// 	})
// }
