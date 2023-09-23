package Guild

import (
	"antegr.al/chatanium-bot/v1/src/Handlers"
	Internal "antegr.al/chatanium-bot/v1/src/Internal"
	"antegr.al/chatanium-bot/v1/src/Log"
	"antegr.al/chatanium-bot/v1/src/Schema"
	"github.com/bwmarrin/discordgo"
)

func Handle(client *discordgo.Session) {
	// TODO: Database (Save Guild ID, etc.)
	var GuildCmds []Commands

	// Register all commands from all guilds
	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		Log.Verbose.Printf("Joined Guild: %v (%v)", g.Name, g.ID)

		AllowedModules := GetModulesByACL(g.ID)

		Guild := Commands{
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

	// Remove all commands from left guilds
	client.AddHandler(func(s *discordgo.Session, g *discordgo.GuildDelete) {
		Log.Verbose.Printf("Left Guild: %v (%v)", g.Name, g.ID)

		for i, v := range GuildCmds {
			if v.GuildID == g.ID {
				v.RemoveSchema()
				GuildCmds = append(GuildCmds[:i], GuildCmds[i+1:]...)
				break
			}
		}
	})

	Internal.MemberLogger(client)
	Internal.MessageLogger(client)
}
