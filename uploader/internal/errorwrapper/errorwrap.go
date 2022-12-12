package errorwrapper

import (
	"commiter/internal/api"
	"commiter/internal/model"
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

type descritionError string

func HandError(e error, extConn *api.ExternalConnection, description ...descritionError) error {
	if cap(description) > 0 {
		sendError(extConn, e.Error(), description[0])
	} else {
		sendError(extConn, e.Error(), "")
	}

	return e
}

func sendError(ec *api.ExternalConnection, errorString string, description descritionError) {

	idChat := IdAdmin(ec.DB)
	msg := tgbotapi.NewMessage(idChat, fmt.Sprint(errorString))
	ec.Bot.Send(msg)
}

func IdAdmin(db *sqlx.DB) int64 {
	us := model.Users{}
	err := db.GetContext(context.TODO(),
		&us,
		"select u.tgid as tgid from users u where u.is_admin = true ")

	if err != nil && !strings.Contains(err.Error(),
		"sql: no rows in result set") {
		return 0
	}
	return us.TgId
}
