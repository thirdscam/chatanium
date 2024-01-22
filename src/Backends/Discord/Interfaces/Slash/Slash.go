package Slash

import (
	"fmt"
	"reflect"

	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/bwmarrin/discordgo"
)

// Manager for slash commands.
// Only one manager can be used for one guild
type CommandManager struct {
	// Commands for serve. (Commands is a map of command schema and handler).
	Commands Commands
	// Discord session from discordgo
	Client *discordgo.Session
	// Guild ID, not Guild Name
	GuildID string
}

func (t *CommandManager) Start() {
	Log.Verbose.Printf("%s > Adding schemas...", t.GuildID)

	// 1. Register commands to discord (only schema. handler is not registered)
	for schema := range t.Commands {
		schema, err := t.Client.ApplicationCommandCreate(t.Client.State.User.ID, t.GuildID, schema)
		if err != nil {
			Log.Error.Fatalf("%s: Cannot create '%v' command: %v", t.GuildID, schema.Name, err)
		}
	}

	Log.Verbose.Printf("%s > Starting CommandManager...", t.GuildID)

	// 2. Handle for slash commands when user input
	t.Client.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// 2-1. Retrieve command name from interaction (input from user)
		reqCmd := i.ApplicationCommandData().Name

		// 2-2. search command name from registered commands
		//      if command name is not found, ignore this event
		for schema, handle := range t.Commands {
			if schema.Name == reqCmd { // schema.Name is command name
				handle(s, i) // call handler
				break
			}
		}
	})
}

// Add command to command manager.
func (t *CommandManager) Add(cmd Commands) {
	Log.Verbose.Printf("%s > Adding commands...", t.GuildID)

	// 1. Add command to discord (register command)
	for schema := range cmd {
		schema, err := t.Client.ApplicationCommandCreate(t.Client.State.User.ID, t.GuildID, schema)
		if err != nil {
			Log.Error.Fatalf("%s: Cannot create '%v' command: %v", t.GuildID, schema.Name, err)
		}
	}

	// 2. Add command to command manager of stored commands
	for schema, handle := range cmd {
		// TODO(Security): thread-safe map (use WaitGroup or Mutex)
		t.Commands[schema] = handle
	}
}

// Remove command from command manager.
func (t *CommandManager) Remove(cmd Commands) {
	Log.Verbose.Printf("%s > Removing commands...", t.GuildID)

	// 1. Remove command from discord (unregister command)
	for schema := range cmd {
		err := t.Client.ApplicationCommandDelete(t.Client.State.User.ID, t.GuildID, schema.ID)
		if err != nil {
			Log.Error.Fatalf("%s: Cannot delete '%v' command: %v", t.GuildID, schema.Name, err)
		}
	}

	// 2. Remove command from command manager of stored commands
	for schema := range cmd {
		delete(t.Commands, schema)
	}
}

// Remove all command from command manager.
func (t *CommandManager) unloadCommand() {
	Log.Verbose.Printf("%s > Unloading commands...", t.GuildID)
	t.Remove(t.Commands)
}

// Vaildate schema of commands.
func (t *CommandManager) vaildate() {
	// Compare from each stored commands and registered commands.
	st, err := t.Client.Guild(t.GuildID)
	if err != nil {
		Log.Error.Fatalf("G:%s > Cannot get guild: %v", t.GuildID, err)
	}

	// 1. Get registered commands from discord
	cmd, err := t.Client.ApplicationCommands(t.Client.State.User.ID, st.ID)
	if err != nil {
		Log.Error.Fatalf("G:%s > Cannot get slash commands: %v", t.GuildID, err)
	}

	// 2. Get stored commands and registered commands.
	remote := make(map[string]string, len(cmd))
	for _, ac := range cmd {
		remote[ac.Name] = fmt.Sprintf("%+v", ac)
	}

	local := make(map[string]string, len(t.Commands))
	for schema := range t.Commands {
		local[schema.Name] = fmt.Sprintf("%+v", schema)
	}

	result := reflect.DeepEqual(remote, local)
	if result == false {
		Log.Warn.Fatalf("[Integrity] G:%s > Stored commands and registered commands are not same. Please check this guild commands.", t.GuildID)
	}
}

// Manage and store CommandManager for each guild.
type Guild struct {
	// Discord session from discordgo
	Client *discordgo.Session
	// Store CommandManager for each guild
	cmdMgrs []CommandManager
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
