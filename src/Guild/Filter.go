package Guild

import "antegr.al/chatanium-bot/v1/src/Log"

func GetModulesByACL(GuildID string) []string {
	switch GuildID {
	case "919823370600742942":
		return []string{"ping", "snowflake2time"}
	default:
		Log.Warn.Printf("%s > No ACL defined, using default", GuildID)
		return []string{}
	}
}
