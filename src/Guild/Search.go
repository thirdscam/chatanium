package Guild

import (
	"antegr.al/chatanium-bot/v1/src/Log"
	"github.com/bwmarrin/discordgo"
)

func SearchUsernameByUID(client *discordgo.Session, uid, gid string) string {
	st, err := client.GuildMember(gid, uid)
	if err != nil {
		Log.Error.Fatalf("Failed to get member: %v", err)
	}
	return st.User.Username
}
