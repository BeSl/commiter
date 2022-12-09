package server

import (
	"commiter/internal/api"
	"commiter/internal/database"
	"commiter/internal/model"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

func createNewJob(w http.ResponseWriter, r *http.Request, ms *MServer) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var upl api.Uplder

	err = json.Unmarshal(b, &upl)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = regMessageDB(&upl, ms)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))

	} else {
		w.WriteHeader(200)
		w.Write([]byte("OK, ready"))
	}
}

func checkedStatusQueues(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	st := StatusQ{}
	err := db.Get(&st, "Select COUNT(*)as CountW FROM commit_tasks WHERE processed=false")

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(200)
		w.Write([]byte(st.CountW))
	}

}

func createTablesDB(db *sqlx.DB) error {

	cq := database.NewCQuery()
	db.MustExec(cq.SchemaUser())
	db.MustExec(cq.SchemaEProc())
	return nil

}

func regMessageDB(upl *api.Uplder, ms *MServer) error {

	us := model.Users{}
	err := ms.db.Get(&us, "Select u.id, u.extid, u.name from users as u where u.extid=$1 LIMIT 1", upl.User.Extid)

	if err != nil && !strings.Contains(err.Error(), "sql: no rows in result set") {
		return err
	}

	tx := ms.db.MustBegin()
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
		return err
	}

	return nil

}
