package errorwrapper

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

func HandError(err error, db *sqlx.DB, bot *tgbotapi.BotAPI, idChat int64) error {

	sendError(db, bot, err.Error(), idChat)

	return err
}

func sendError(db *sqlx.DB, bot *tgbotapi.BotAPI, errorString string, chat int64) {

	msg := tgbotapi.NewMessage(chat, fmt.Sprint(errorString))
	bot.Send(msg)
}
