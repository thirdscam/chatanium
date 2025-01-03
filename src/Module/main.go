package module

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"plugin"

	"antegr.al/chatanium-bot/v1/src/Util/Log"
)

// ModuleManager is managed on a per-backend basis.
// Modules that match the Identifier specified by the Backend are assigned to each Backend.
//
// If a Backend named Discord has an Identifier of "discord",
// then modules based on the discord backend should have a package name of "chatanium_discord".
type ModuleManager struct {
	Identifier string
	Modules    map[string]func()
}

func (t *ModuleManager) Start() {
	t.Modules = make(map[string]func())

	Log.Verbose.Printf("Module@%s > Loading modules...", t.Identifier)

	files, err := os.ReadDir("./modules")
	if err != nil {
		log.Fatal(err)
	}

	Log.Verbose.Printf("Module@%s > Found %d files.", t.Identifier, len(files))

	// Load all plugins
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".so" {
			plug, err := plugin.Open("./modules/" + file.Name())
			if err != nil {
				Log.Warn.Printf("Module@%s > Unknown (%s) > Failed to load plugin: %v", t.Identifier, file.Name(), err)
				continue
			}

			// 1. Build module
			module := Module{
				Backend: t.Identifier,
			}

			ok := module.Build(file.Name(), plug)
			if !ok {
				Log.Warn.Printf("Module@%s > %s (%s) > Failed to get module info from plugin.", t.Identifier, module.Name, file.Name())
				continue
			}

			// 2. Check plugin integrity
			f, err := os.Open("./modules/" + file.Name())
			if err != nil {
				Log.Warn.Printf("Module@%s > %s (%s) > Failed to open plugin: %v", t.Identifier, module.Name, file.Name(), err)
				continue
			}

			h := sha256.New()
			if _, err := io.Copy(h, f); err != nil {
				Log.Warn.Printf("Module@%s > %s (%s) > Failed to calculate hash of plugin: %v", t.Identifier, module.Name, file.Name(), err)
				continue
			}
			hash := hex.EncodeToString(h.Sum(nil))

			Log.Info.Printf("Module@%s > %s (%s) > Inserted. (%s)", t.Identifier, module.Name, file.Name(), hash)
		}
	}

	// Start all plugins
	for name, startFunc := range t.Modules {
		Log.Info.Printf("Starting plugin: %s...", name)
		startFunc()
	}
}

func (t *ModuleManager) Shutdown() {
}
