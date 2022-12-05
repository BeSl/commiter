package api

import (
	"commiter/internal/model"
)

type Uplder struct {
	User        *model.Users             `json:"user"`
	ExtDataProc *model.ExtDataProcessors `json:"extprc"`
}

func NewUplder(u *model.Users, edp *model.ExtDataProcessors) *Uplder {
	return &Uplder{
		User:        u,
		ExtDataProc: edp,
	}
}
