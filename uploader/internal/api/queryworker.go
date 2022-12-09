package api

import (
	"commiter/internal/model"
	"time"
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


