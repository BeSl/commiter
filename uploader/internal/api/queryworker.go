package api

import (
	"commiter/internal/model"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

type Uplder struct {
	User           *model.Users             `json:"author"`
	DataProccessor *model.ExtDataProcessors `json:"DataProccessor"`
	TextCommit     string                   `json:"textCommit"`
	Dataevent      string                   `json:"dataevent"`
}

func NewUplder(u *model.Users, edp *model.ExtDataProcessors) *Uplder {
	return &Uplder{
		User:           u,
		DataProccessor: edp,
	}
}

type CurrentJob struct {
	Id            int64     `db:"id"`
	Name          string    `db:"j_name"`
	DateJob       time.Time `db:"j_date"`
	ErrorDescript string    `db:"j_error"`
	Prim          string    `db:"j_prim"`
}

func NewCurrentJob() *CurrentJob {
	return &CurrentJob{}
}

type ExternalConnection struct {
	DB  *sqlx.DB
	Bot *tgbotapi.BotAPI
}
