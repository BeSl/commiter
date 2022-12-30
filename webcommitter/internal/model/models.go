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

	Authorcommit struct {
		gorm.Model
		Name        string `json:"name" `
		ExtRef      string `json:"extId"`
		GitEmail    string `json:"email"`
		IsAdmin     bool
		TgChannelID int64
	}

	DataProccessor struct {
		gorm.Model
		ExtRef     string    `json:"extID"`
		Name       string    `json:"name" `
		Dateevent  time.Time `json:"dateevent"`
		Base64data string    `json:"base64data"`
		TypeData   string    `json:"type"`
	}

	Commit struct {
		gorm.Model
		AuthorCommitID int
		AuthorCommit   Authorcommit `json:"author" gorm:"foreignKey:AuthorCommitID"`
		ProccessorID   int
		Proccessor     DataProccessor `json:"DataProccessor" gorm:"foreignKey:ProccessorID"`
		Data           string         `json:"DataProccessor.base64data"`
		TextCommit     string         `json:"textCommit" `
		Source         string         `json:"source"`
		ItsDone        bool
	}

	ResponseData struct {
		AuthorEmail   string
		Branch        string
		Data          string
		NameProcessor string
		TextCommit    string
		IDCommit      int
	}
)
