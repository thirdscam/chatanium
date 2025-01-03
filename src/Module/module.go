package module

import (
	"plugin"

	"antegr.al/chatanium-bot/v1/src/Util/Log"
)

type Module struct {
	ManifestVersion int

	Name       string
	Backend    string
	Version    string
	Author     string
	Repository string

	entryPoint func()
}

// Getting module info
func (t *Module) Build(filename string, plugin *plugin.Plugin) bool {
	t.Name = "Unknown"

	// 1. Check MANIFEST_VERSION variable from plugin
	manifestSymbol, err := plugin.Lookup("MANIFEST_VERSION")
	if err != nil {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin does not export MANIFEST_VERSION.", t.Backend, t.Name, filename)
		return false
	}

	ManifestVersion, ok := manifestSymbol.(*int)
	if !ok {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin has invalid MANIFEST_VERSION type.", t.Backend, t.Name, filename)
		return false
	}

	t.ManifestVersion = *ManifestVersion

	// 2. Check NAME variable from plugin
	nameSymbol, err := plugin.Lookup("NAME")
	if err != nil {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin does not export NAME: %v", t.Backend, t.Name, filename, err)
		return false
	}

	name, ok := nameSymbol.(*string)
	if !ok {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin has invalid NAME type.", t.Backend, t.Name, filename)
		return false
	}

	t.Name = *name

	// 3. Check BACKEND variable from plugin
	backendSymbol, err := plugin.Lookup("BACKEND")
	if err != nil {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin does not export BACKEND: %v", t.Backend, t.Name, filename, err)
		return false
	}

	backend, ok := backendSymbol.(*string)
	if !ok {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin has invalid BACKEND type.", t.Backend, t.Name, filename)
		return false
	}

	if *backend != t.Backend {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin is not for this backend (%s != %s).", t.Backend, t.Name, filename, *backend, t.Backend)
		return false
	}

	// 4. Check VERSION variable from plugin
	versionSymbol, err := plugin.Lookup("VERSION")
	if err != nil {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin does not export VERSION: %v", t.Backend, t.Name, filename, err)
		return false
	}

	version, ok := versionSymbol.(*string)
	if !ok {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin has invalid VERSION type.", t.Backend, t.Name, filename)
		return false
	}

	t.Version = *version

	// 5. Check AUTHOR variable from plugin
	authorSymbol, err := plugin.Lookup("AUTHOR")
	if err != nil {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin does not export AUTHOR: %v", t.Backend, t.Name, filename, err)
		return false
	}

	author, ok := authorSymbol.(*string)
	if !ok {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin has invalid AUTHOR type.", t.Backend, t.Name, filename)
		return false
	}

	t.Author = *author

	// 6. Check REPOSITORY variable from plugin
	repositorySymbol, err := plugin.Lookup("REPOSITORY")
	if err != nil {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin does not export REPOSITORY: %v", t.Backend, t.Name, filename, err)
		return false
	}

	repository, ok := repositorySymbol.(*string)
	if !ok {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin has invalid REPOSITORY type.", t.Backend, t.Name, filename)
		return false
	}

	t.Repository = *repository

	// 7. Check Entry point from plugin
	startSymbol, err := plugin.Lookup("Start")
	if err != nil {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin does not export Start function: %v", t.Backend, t.Name, filename, err)
		return false
	}

	startFunc, ok := startSymbol.(func())
	if !ok {
		Log.Warn.Printf("Module@%s > %s (%s) > Plugin has invalid Start function type.", t.Backend, t.Name, filename)
		return false
	}

	t.entryPoint = startFunc

	return true
}
