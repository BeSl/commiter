package storage

import (
	"commiter/internal/model"
	"database/sql"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	//"commiter/internal/model"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

type CommitStatus struct {
	CommitCount string
}

type Storage struct {
	DB  *sqlx.DB
	Bot *tgbotapi.BotAPI
}

func NewStorage(db *sqlx.DB, bot *tgbotapi.BotAPI) *Storage {
	return &Storage{
		DB:  db,
		Bot: bot,
	}
}

func (s *Storage) AddNewRequest(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var upl model.DataCommit

	err = json.Unmarshal(b, &upl)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = s.regMessageDB(&upl)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))

	} else {
		w.WriteHeader(200)
		w.Write([]byte("OK, ready"))
	}
}

func (s *Storage) CheckedStatusQueues(w http.ResponseWriter, r *http.Request) {
	st := CommitStatus{}
	err := s.DB.Get(&st, "Select COUNT(*)as CommitCount FROM commit_tasks WHERE processed=false")

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(200)
		w.Write([]byte(st.CommitCount))
	}

}

func (s *Storage) CreateTablesDB() error {

	// cq := database.NewCQuery()
	// s.ExtConn.DB.MustExec(cq.SchemaUser())
	// s.ExtConn.DB.MustExec(cq.SchemaEProc())
	return nil

}

func (s *Storage) regMessageDB(upl *model.DataCommit) error {

	us := model.User{}
	err := s.DB.GetContext(context.TODO(),
		&us,
		"Select u.id, u.extid, u.name from users as u where u.extid=$1 LIMIT 1",
		upl.User.ExtID)

	if err != nil && !strings.Contains(err.Error(),
		"sql: no rows in result set") {
		//notUsers := fmt.Sprintf("Not users %s id=%s", upl.User.Name, upl.User.Extid)
		return err
	}

	tx := s.DB.MustBegin()
	if us.Base.ID == 0 {
		tx.MustExec("INSERT INTO users (extId, name) VALUES ($1, $2)",
			string(upl.User.ExtID), string(upl.User.FullName))
	}

	tx.MustExec("INSERT INTO commit_tasks (extId,name,base64data,type,userid,textcommit) VALUES ($1, $2, $3, $4, $5, $6)",
		upl.DataProccessor.ExtID,
		upl.DataProccessor.Name,
		upl.DataProccessor.Base64data,
		upl.DataProccessor.Type,
		us.Base.ID,
		upl.TextCommit)

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) FindUserByID(id string) (*model.User, error) {

	u := &model.User{}
	q := "Select id,extID,FullName,GitEmail,IsAdmin,TGGit FROM user WHERE id = $1"
	err := s.runSelectQuery(q, &u, u.Base.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Storage) FindAdmin() (*model.User, error) {
	return nil, nil
}

func (s *Storage) runSelectQuery(query string, nO interface{}, args ...interface{}) error {

	err := s.DB.QueryRowxContext(context.TODO(),
		query, args[0:]...).Scan()

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) FindLastCommit() (*model.DataCommit, error) {

	q :=
		`SELECT 
			coalesce(u.gitlogin, '') as gitlogin,
			u.name as username,
			ct.name as name,
			ct.id as id,
			ct.type as type,
			ct.base64data as base64data,
			coalesce(ct.textcommit, 'not text') as commit 
		FROM 
			commit_tasks ct 
				left join users u 
				on u.id= ct.userid 
		WHERE 
		 		ct.processed =false 
		ORDER BY ct.id 
		LIMIT 1`

	dw := &model.DataCommit{}
	err := s.DB.Get(dw, q)

	if len(dw.User.GitEmail) == 0 {
		return nil, errors.New("Не заполнен пользователь " + dw.User.FullName)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return dw, nil
}
