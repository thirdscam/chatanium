/*
Copyright © 2025 ANTEGRAL <antegral@antegral.net>
*/
package cmd

import (
	"github.com/spf13/cobra"
	discord "github.com/thirdscam/chatanium/src/Backends/Discord"
	"github.com/thirdscam/chatanium/src/Database"
	"github.com/thirdscam/chatanium/src/Util"
	"github.com/thirdscam/chatanium/src/Util/Log"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the modules",
	Long:  `Start the modules`,
	Run: func(cmd *cobra.Command, args []string) {
		// Init the logging system
		Log.Init()

		// Init the environment variables
		Util.InitEnv()

		Log.Info.Println("Antegral/Chatanium: Dynamic Bot Management System")
		Log.Info.Println("Press CTRL+C to shutdown.")

		// Ignite Database
		database := &Database.DB{}
		database.Start()
		defer database.Shutdown()

		// Create a new Discord session using the provided bot token.
		discord := &discord.Backend{
			Token:    Util.GetEnv("DISCORD_TOKEN"),
			Database: database,
		}
		discord.Start()
		defer discord.Shutdown()

		// Wait for a signal to shutdown
		Util.WaitSignal()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
