package server

import (
	"commiter/internal/api"
	"commiter/internal/model"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func (ms *MServer) uploadtoquery(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/uploadtoquery" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Error type Method", 405)
		return
	}

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

func (ms *MServer) pingService(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ping" {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("pong"))
}

func (ms *MServer) CreateTables(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/crtab" {
		http.NotFound(w, r)
		return
	}
	ms.db.MustExec(schema)
	ms.db.MustExec(schemaUser)
	ms.db.MustExec(schemaEProc)
	w.WriteHeader(200)

}

var schema = `
CREATE TABLE person (
    first_name text,
    last_name text,
    email text
);`

var schemaUser = `
CREATE TABLE users (
    id serial PRIMARY KEY,
	extId UUID, 
	name text,
    is_admin boolean,
    gitlogin text,
	tgID text
);`

var schemaEProc = `
CREATE TABLE extprocVersion (
    id serial PRIMARY KEY,
	authorversion UUID,
	extId UUID, 
	name text,
    BinaryData text,
    Filename text,
	Processed boolean
);`

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
	tx.Commit()

	return nil

}
