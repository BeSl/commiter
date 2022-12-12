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

func HandError(err error, extConn *api.ExternalConnection, description string) error {
	if len(description) > 0 {
		sendError(extConn, err.Error(), description)
	} else {
		sendError(extConn, err.Error(), "not")
	}

	return err
}

func sendError(ec *api.ExternalConnection, errorString string, description string) {

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
