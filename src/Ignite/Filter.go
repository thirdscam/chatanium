package Ignite

import "antegr.al/chatanium-bot/v1/src/Log"

func GetModulesByACL(GuildID string) []string {
	// TODO: Get ACL from database
	switch GuildID {
	case "919823370600742942":
		return []string{"ping", "echo"}
	default:
		Log.Warn.Printf("%s > No ACL defined, using default", GuildID)
		return []string{}
	}
}
