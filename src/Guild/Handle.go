package Guild

import (
	"antegr.al/chatanium-bot/v1/src/Database"
	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	"antegr.al/chatanium-bot/v1/src/Handlers"
	Internal "antegr.al/chatanium-bot/v1/src/Internal"
	"antegr.al/chatanium-bot/v1/src/Log"
	"antegr.al/chatanium-bot/v1/src/Schema"
	"github.com/bwmarrin/discordgo"
)

func Handle(client *discordgo.Session, db *db.PrismaClient) {
	var GuildCmdStorage []Commands

	// Register all commands from all guilds
	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		Log.Verbose.Printf("Joined Guild: %v (%v)", g.Name, g.ID)

		AllowedModules := GetModulesByACL(g.ID)

		GuildCmd := Commands{
			Schema:   Schema.GetAllowedOnly(AllowedModules),
			Handlers: Handlers.GetAllowedOnly(AllowedModules),
			Client:   client,
			GuildID:  g.ID,
		}

		// Register commands from guild
		GuildCmd.RegisterHandlers()
		GuildCmd.RegisterSchema()

		// Register guild in database
		Database.RegisterGuild(client, db, g)

		GuildCmdStorage = append(GuildCmdStorage, GuildCmd)
	})

	// Remove all commands from left guilds
	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildDelete) {
		Log.Verbose.Printf("Left Guild: %v (%v)", g.Name, g.ID)

		for i, v := range GuildCmdStorage {
			if v.GuildID == g.ID {
				v.RemoveSchema()
				GuildCmdStorage = append(GuildCmdStorage[:i], GuildCmdStorage[i+1:]...)
				break
			}
		}
	})

	// Internal Modules
	Internal.MemberLogger(client, db)
	Internal.MessageLogger(client, db)
}
