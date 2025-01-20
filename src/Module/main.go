package module

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"plugin"

	"github.com/thirdscam/chatanium/src/Util/Log"
)

// ModuleManager is managed on a per-backend basis.
// Modules that match the Identifier specified by the Backend are assigned to each Backend.
//
// If a Backend named Discord has an Identifier of "discord",
// then modules based on the discord backend should have a package name of "chatanium_discord".
type ModuleManager struct {
	Identifier string
	Modules    map[string]Module
}

func (t *ModuleManager) Load() {
	t.Modules = make(map[string]Module)
	Log.Verbose.Printf("Module@%s > Loading modules...", t.Identifier)

	files, err := os.ReadDir("./modules")
	if err != nil {
		log.Fatal(err)
	}

	Log.Verbose.Printf("Module@%s > Found %d files.", t.Identifier, len(files))

	// Load all plugins
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".so" {
			continue
		}

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
		t.Modules[hash] = module
		Log.Info.Printf("Module@%s > %s (%s) > Inserted. (%s)", t.Identifier, module.Name, file.Name(), hash[:8])

	}
	Log.Verbose.Printf("Module@%s > Loaded %d plugins.", t.Identifier, len(t.Modules))
}

// Start all plugins
func (t *ModuleManager) Start() {
	Log.Info.Printf("Module@%s > Starting %d plugins...", t.Identifier, len(t.Modules))
	for hash, module := range t.Modules {
		go func(m Module, h string) {
			Log.Verbose.Printf("Module@%s > Starting plugin: %s (%s)...", t.Identifier, m.Name, h[:8])
			m.entryPoint()
		}(module, hash)
	}
}

func (t *ModuleManager) Shutdown() {
}
