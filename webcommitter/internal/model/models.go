package model

import (
	"time"

	"gorm.io/gorm"
)

type ExtDataQuery struct {
	User           *Users             `json:"author"`
	DataProccessor *ExtDataProcessors `json:"DataProccessor"`
	TextCommit     string             `json:"textCommit"`
	Dataevent      string             `json:"dataevent"`
}

type (
	DataCommit struct {
		User           *User              `json:"author"`
		DataProccessor *ExtDataProcessors `json:"DataProccessor"`
		TextCommit     string             `json:"textCommit"`
		Dataevent      string             `json:"dataevent"`
	}

	User struct {
		ID       int64  `db:"id"`
		Name     string `db:"name"`
		Deletion bool   `db:"deletion"`
		ExtID    string `db:"extid"`
		FullName string `db:"name"`
		GitEmail string `db:"gitlogin"`
		IsAdmin  bool   `db:"is_admin"`
		TGid     int64  `db:"tgid"`
	}

	Project struct {
		ID          int64  `db:"id"`
		Name        string `db:"name"`
		Deletion    bool   `db:"deletion"`
		IsBlock     bool
		Description string
		GitURL      string
		ProdBranch  string
		DevBranch   string
	}
	DataWork struct {
		Base64data string `db:"base64data"`
		Name       string `db:"name"`
		ID         int64  `db:"id"`
		TypeProc   string `db:"type"`
		UserName   string `db:"username"`
		GitLogin   string `db:"gitlogin"`
		Commit     string `db:"commit"`
	}

	AuthorCommit struct {
		gorm.Model
		Name  string `json:"name"`
		ExtID string `json:"extId"`
	}

	DataProccessor struct {
		gorm.Model
		ExtID      string    `json:"extID"`
		Name       string    `json:"name"`
		Dateevent  time.Time `json:"dateevent"`
		Base64data string    `json:"base64data"`
		TypeData   string    `json:"Отчет"`
	}

	NewCommit struct {
		gorm.Model
		Author         *AuthorCommit   `json:"author"`
		Dataproccessor *DataProccessor `json:"DataProccessor"`
		TextCommit     string          `json:"textCommit"`
		Source         string          `json:"source"`
	}
)
