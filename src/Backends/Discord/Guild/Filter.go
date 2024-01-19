package Guild

import "antegr.al/chatanium-bot/v1/src/Util/Log"

func GetModuleStringByACL(GuildID string) []string {
	switch GuildID {
	case "919823370600742942":
		return []string{"ping", "snowflake2time"}
	default:
		Log.Warn.Printf("%s > Undefined ACL for Guild. using default", GuildID)
		return []string{}
	}
}
