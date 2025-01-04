package Slash

import (
	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/bwmarrin/discordgo"
)

// Handle for slash commands when user input
func (t *CommandManager) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 1. Retrieve command name from interaction (input from user)
	reqCmd := i.ApplicationCommandData().Name

	// 2. search command name from registered commands
	//    if command name is not found, ignore this event
	for schema, handle := range t.Commands {
		if schema.Name == reqCmd { // schema.Name is command name
			handle(s, i) // call handler
			break
		}
	}
}

// Manage and store CommandManager for each guild.
type Guild struct {
	// Discord session from discordgo
	Client *discordgo.Session
	// Store each CommandManager of guild
	cmdMgrs []CommandManager
}

// Start slash command interface.
func (t *Guild) Start() {
	// Add handler for slash command event (InteractionCreate)
	t.Client.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// 1. Get guild ID from interaction
		guildId := i.GuildID
		isFindCmdMgr := false

		// 2. Search CommandManager for guild
		for _, v := range t.cmdMgrs {
			if v.GuildID == guildId {
				// 2-1. Execute command: but if command is not found, ignore this event
				isFindCmdMgr = true
				v.HandleCommand(s, i)
				break
			}
		}

		if !isFindCmdMgr {
			Log.Warn.Printf("[Integrity] G:%s > Cannot find CommandManager for this guild. Please check this guild.", guildId)
		}
	})
}

// Register slash commands to guild.
// This function used on guild create event.
func (t *Guild) OnGuildCreated(id string) {
	// TODO(Feature): ACL of slash commands for each guild

	// 1. Make CommandManager for guild
	cmdMgrs := CommandManager{
		Client:   t.Client,
		Commands: getCommands(id),
		GuildID:  id,
	}

	// 2. Start CommandManager
	cmdMgrs.Start()
	cmdMgrs.Vaildate() // Validate commands

	// 3. Append to CommandManager storage
	t.cmdMgrs = append(t.cmdMgrs, cmdMgrs)
}

// Remove slash commands from guild.
// This function used on guild delete event.
func (t *Guild) OnGuildDeleted(id string) {
	for i, v := range t.cmdMgrs {
		if v.GuildID == id {
			// FIXME(Performance): remove CommandManager for lefted guild nevertheless handler goroutine is may running,
			// so this goroutine maybe not stopped. The current way to fix this is to periodically restart
			// the server to remove handlers for guilds that have lefted.
			t.cmdMgrs = append(t.cmdMgrs[:i], t.cmdMgrs[i+1:]...)
			break
		}
	}
}
