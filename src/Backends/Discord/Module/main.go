package Module

import (
	"github.com/thirdscam/chatanium/src/Backends/Discord/Interface/Slash"
	module "github.com/thirdscam/chatanium/src/Module"
	"github.com/thirdscam/chatanium/src/Util/Log"
)

type DiscordModule struct {
	module.Module
	Commands Slash.Commands
}

func ConvertDiscordModule(m module.Module) DiscordModule {
	plugin := m.RawPlugin

	slashCmdSymbol, err := plugin.Lookup("DEFINE_SLASHCMD")
	if err != nil {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin does not export DEFINE_SLASHCMD: %v", m.Backend, m.Name, m.Filename, err)
	}

	slashCmd, ok := slashCmdSymbol.(*Slash.Commands)
	if !ok {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin has invalid DEFINE_SLASHCMD type.", m.Backend, m.Name, m.Filename)
	}

	if slashCmd == nil {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin has no defined slash commands.", m.Backend, m.Name, m.Filename)
		return DiscordModule{
			Module:   m,
			Commands: Slash.Commands{},
		}
	}

	return DiscordModule{
		Module:   m,
		Commands: *slashCmd,
	}
}

// func GetCommandMap(modules map[string]DiscordModule) map[string]Slash.Commands {
// 	for k, v := range modules {
// 	}
// }
