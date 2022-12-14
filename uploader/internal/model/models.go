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

	BaseValue struct {
		ID       int64  `db:"id"`
		Name     string `db:"name"`
		Deletion bool   `db:"deletion"`
	}

	User struct {
		Base     *BaseValue
		ExtID    string
		FullName string
		GitEmail string
		IsAdmin  bool
		TGid     int64
	}

	Project struct {
		Base        *BaseValue
		IsBlock     bool
		Description string
		GitURL      string
		ProdBranch  string
		DevBranch   string
	}
)
