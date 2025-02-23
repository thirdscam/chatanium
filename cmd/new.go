/*
Copyright © 2025 ANTEGRAL <antegral@antegral.net>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/thirdscam/chatanium/cmd/tui/newcmdinput"
	"github.com/thirdscam/chatanium/src/Util/Log"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new chatanium module",
	Long: `This command creates a new chatanium module.
it will create a new module with the given name.`,
	Run: func(cmd *cobra.Command, args []string) {
		Log.Init()

		// Getting module name from user
		model, err := tea.NewProgram(newcmdinput.GetModel()).Run()
		if err != nil {
			Log.Error.Fatalf("Error running TUI model: %v", err)
		}

		m, ok := model.(newcmdinput.Model)
		if !ok {
			Log.Error.Fatalf("Error getting value from tui program!")
		}

		backendName := "discord" // TODO: user can choose backend (maybe need to add more backends)
		moduleName := m.GetValue()
		if moduleName == "" {
			Log.Error.Fatalf("Module name cannot be empty")
		}

		// Getting Username
		user, err := user.Current()
		if err != nil {
			Log.Error.Fatalf("Failed to get current user: %v", err)
		}

		userName := user.Username
		if userName == "" {
			Log.Error.Fatalf("Failed to get current user: username is empty")
		}

		isVaildArgs := isVaildBackend(backendName)
		if !isVaildArgs {
			Log.Error.Fatalf("Invalid backend name: %s", backendName)
		}

		Log.Info.Printf("Creating new module: %s", moduleName)
		Log.Info.Printf("Backend: %s", backendName)

		// Get go version
		GO_VERSION := runtime.Version()
		if GO_VERSION == "" {
			Log.Error.Fatalf("Failed to get go version")
		}

		// Make module path
		modulePath, _ := filepath.Abs(filepath.Join("modules", moduleName))

		// Create module directory
		if err := os.MkdirAll(modulePath, 0o755); err != nil {
			Log.Error.Fatalf("Failed to create modules directory: %v", err)
		}

		// Check go compiler binary
		goBinPath := filepath.Join(os.Getenv("GOPATH"), "bin", GO_VERSION)
		if _, err := os.Stat(goBinPath); os.IsNotExist(err) {
			Log.Error.Fatalf("Failed to find go compiler binary: %v", err)
		}

		Log.Info.Printf("Found go compiler binary: %s", goBinPath)

		// Create go.mod
		execCmd := exec.Command(goBinPath, "mod", "init", fmt.Sprintf("github.com/%s/%s", userName, moduleName))
		execCmd.Dir = modulePath
		if err := execCmd.Run(); err != nil {
			Log.Error.Fatalf("Failed to create go.mod: %v", err)
		}

		// append runtime dependencies
		appendModuleDeps := "require github.com/thirdscam/chatanium v1.0.0-local\n\nreplace github.com/thirdscam/chatanium v1.0.0-local => ./../.."
		f, err := os.OpenFile(filepath.Join(modulePath, "go.mod"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
		if err != nil {
			Log.Error.Fatalf("Failed to open go.mod: %v", err)
		}

		if _, err := f.WriteString(appendModuleDeps); err != nil {
			Log.Error.Fatalf("Failed to append module dependencies: %v", err)
		}

		f.Close()
		Log.Info.Printf("Appended module dependencies to go.mod")

		// Create main.go with constants
		if err := os.WriteFile(filepath.Join(modulePath, "main.go"), []byte(getMainGoCode(backendName, moduleName, userName)), 0o644); err != nil {
			Log.Error.Fatalf("Failed to create main.go: %v", err)
		}

		// Run go mod tidy
		execCmd = exec.Command(goBinPath, "mod", "tidy")
		execCmd.Dir = modulePath
		if err := execCmd.Run(); err != nil {
			Log.Error.Fatalf("Failed to run go mod tidy: %v", err)
		}
		Log.Info.Printf("go mod tidy Completed")

		Log.Info.Printf("Successfully created module at %s!", modulePath)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func isVaildBackend(backendName string) bool {
	return backendName == "discord"
}

func getMainGoCode(backendName string, moduleName string, username string) string {
	return fmt.Sprintf(`package main

import (
	"github.com/thirdscam/chatanium/src/Util/Log"
)

const MANIFEST_VERSION = 1

const (
	NAME       = "%s"
	BACKEND    = "%s"
	VERSION    = "0.0.1"
	AUTHOR     = "%s"
	REPOSITORY = "github:%s/%s"
)

func Start() {
	// Add your code here!
	Log.Info.Println("Hello World! I'm %s!")
}
`, moduleName, backendName, username, username, moduleName, moduleName)
}
