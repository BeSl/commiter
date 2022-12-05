package server

import (
	"commiter/internal/api"
	"commiter/internal/model"
	"encoding/json"
	"io/ioutil"
	"net/http"
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

	w.WriteHeader(200)
	w.Write([]byte("OK, ready"))

}

func (ms *MServer) pingService(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ping" {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("pong"))
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

	// ms.db.MustExec(schema)
	// ms.db.MustExec(schemaUser)
	//ms.db.MustExec(schemaEProc)

	tx := ms.db.MustBegin()
	tx.MustExec("INSERT INTO users (extId, name) VALUES ($1, $2)", upl.User.Id, upl.User.Name)
	tx.MustExec("INSERT INTO extprocVersion (extId, name, BinaryData,Filename) VALUES ($1, $2, $3, $4)", upl.ExtDataProc.UidBase, upl.ExtDataProc.Name, upl.ExtDataProc.BinaryData, upl.ExtDataProc.Filename)

	us := model.NewUsers()

	err := ms.db.Get(&us, "Select id from users where extid=$1", upl.User.Id)

	if err != nil {
		return err
	}
	tx.Commit()

	return nil

}
