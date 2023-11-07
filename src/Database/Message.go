package Database

import (
	"context"
	"errors"

	db "antegr.al/chatanium-bot/v1/src/Database/Internal"
	util "antegr.al/chatanium-bot/v1/src/Util"
	"antegr.al/chatanium-bot/v1/src/Util/Log"
)

func GetMessageInfo(gid, mid string, database *db.PrismaClient) *db.MessagesModel {
	msg, err := database.Messages.FindUnique(
		db.Messages.MessageID.Equals(util.StringToBigint(mid)),
	).With(
		db.Messages.Users.Fetch(),
	).Exec(
		context.Background(),
	)
	if errors.Is(err, db.ErrNotFound) {
		Log.Error.Printf("G:%s | M:%s > Cannot find message : %v", gid, mid, err)
		return nil
	} else if err != nil {
		Log.Error.Printf("G:%s | M:%s > Failed to find message : %v", gid, mid, err)
		return nil
	}

	return msg
}
