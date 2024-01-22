package Slash

import (
	"github.com/bwmarrin/discordgo"
)

type Commands map[*discordgo.ApplicationCommand]func(s *discordgo.Session, i *discordgo.InteractionCreate)

// get commands for guild from pre-defined commands and modules.
func getCommands(guildId string) Commands {
	// TODO(Security): Support ACL
	// TODO(Feature): load commands from third-party modules
	return PREDEFINED_COMMANDS
}

// func GetAllHandler() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
// 		// Built-in commands configuration
// 		// Write your handler of COMMANDS for module here
// 		// Example:
// 		// "command_name": HandleCommand,
// 	}
// }

// func GetOnlyHandler(AllowedModules []string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 	result := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
// 	AllHandlers := GetAllHandler()

// 	// Insert only allowed modules to result
// 	for _, v := range AllowedModules {
// 		for name, fn := range AllHandlers {
// 			if name == v {
// 				result[v] = fn
// 				break
// 			}
// 		}
// 	}

// 	return result
// }

// func GetAllSchema() []*discordgo.ApplicationCommand {
// 	var result []*discordgo.ApplicationCommand
// 	schemas := getDefinedSchema()

// 	for _, v := range schemas {
// 		result = append(result, v())
// 	}

// 	return result
// }

// func GetOnlySchema(AllowedModules []string) []*discordgo.ApplicationCommand {
// 	var result []*discordgo.ApplicationCommand
// 	AllSchemas := getDefinedSchema()

// 	for _, v := range AllowedModules {
// 		for name, schema := range AllSchemas {
// 			if name == v {
// 				result = append(result, schema())
// 				break
// 			}
// 		}
// 	}

// 	return result
// }
