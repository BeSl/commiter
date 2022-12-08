package api

import (
	"commiter/internal/model"
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
