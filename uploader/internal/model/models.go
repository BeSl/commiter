package model

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
)
