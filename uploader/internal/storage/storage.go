package storage

import (
	"commiter/internal/api"
	"commiter/internal/database"
	"commiter/internal/errorwrapper"
	"commiter/internal/model"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

type CommitStatus struct {
	CommitCount string
}

type Storage struct {
	ExtConn *api.ExternalConnection
}

func NewStorage(db *sqlx.DB, bot *tgbotapi.BotAPI) *Storage {
	return &Storage{
		ExtConn: &api.ExternalConnection{
			DB:  db,
			Bot: bot,
		},
	}
}
func (s *Storage) AddNewRequest(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var upl api.Uplder

	err = json.Unmarshal(b, &upl)
	if err != nil {
		http.Error(w,
			errorwrapper.HandError(err, s.ExtConn, "AddNewRequest").
				Error(),
			500)

		return
	}

	err = s.regMessageDB(&upl)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(errorwrapper.HandError(err, s.ExtConn).Error()))

	} else {
		w.WriteHeader(200)
		w.Write([]byte("OK, ready"))
	}
}

func (s *Storage) CheckedStatusQueues(w http.ResponseWriter, r *http.Request) {
	st := CommitStatus{}
	err := s.ExtConn.DB.Get(&st, "Select COUNT(*)as CommitCount FROM commit_tasks WHERE processed=false")

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(errorwrapper.HandError(err, s.ExtConn).Error()))
	} else {
		w.WriteHeader(200)
		w.Write([]byte(st.CommitCount))
	}

}

func (s *Storage) CreateTablesDB() error {

	cq := database.NewCQuery()
	s.ExtConn.DB.MustExec(cq.SchemaUser())
	s.ExtConn.DB.MustExec(cq.SchemaEProc())
	return nil

}

func (s *Storage) regMessageDB(upl *api.Uplder) error {

	us := model.Users{}
	err := s.ExtConn.DB.GetContext(context.TODO(),
		&us,
		"Select u.id, u.extid, u.name from users as u where u.extid=$1 LIMIT 1",
		upl.User.Extid)

	if err != nil && !strings.Contains(err.Error(),
		"sql: no rows in result set") {
		return errorwrapper.HandError(err, s.ExtConn)
	}

	tx := s.ExtConn.DB.MustBegin()
	if us.Id == 0 {
		tx.MustExec("INSERT INTO users (extId, name) VALUES ($1, $2)",
			string(upl.User.Extid), string(upl.User.Name))
	}

	tx.MustExec("INSERT INTO commit_tasks (extId,name,base64data,type,userid) VALUES ($1, $2, $3, $4, $5)",
		upl.DataProccessor.ExtID,
		upl.DataProccessor.Name,
		upl.DataProccessor.Base64data,
		upl.DataProccessor.Type,
		us.Id)

	err = tx.Commit()

	if err != nil {
		return errorwrapper.HandError(err, s.ExtConn)
	}

	return nil
}
