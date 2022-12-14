package comittworker

import (
	"commiter/internal/errorwrapper"
	"commiter/internal/model"
	"commiter/internal/storage"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

const TimeSleepMinute = 1

type CommitCreator struct {
	DB  *sqlx.DB
	Bot *tgbotapi.BotAPI
}

func NewCommitCreator(db *sqlx.DB, bot *tgbotapi.BotAPI) *CommitCreator {
	return &CommitCreator{
		DB:  db,
		Bot: bot,
	}
}

func (cc *CommitCreator) ListenNewTasks() error {

	st := storage.NewStorage(cc.DB, cc.Bot)
	adminUser, err := st.FindAdmin()
	if err != nil {
		return err
	}

	for {
		dataCommit, err := st.FindLastCommit()
		if err != nil {
			errorwrapper.HandError(err, cc.DB, cc.Bot, adminUser.TGid)
		}

		err = createCommit(dataCommit, st)
		if err != nil {
			errorwrapper.HandError(err, cc.DB, cc.Bot, adminUser.TGid)
		}

		time.Sleep(time.Minute * time.Duration(TimeSleepMinute))
	}

}

func createCommit(dc *model.DataCommit, st *storage.Storage) error {

	return nil
}
