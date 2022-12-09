package server

import (
	"net/http"
)

type StatusQ struct {
	CountW string
}

func (ms *MServer) uploadtoquery(w http.ResponseWriter, r *http.Request) {

	if chekPathRequest(w, r, "/uploadtoquery") == false {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Error type Method", 405)
		return
	}

	createNewJob(w, r, ms)

}

func (ms *MServer) StatusQueue(w http.ResponseWriter, r *http.Request) {

	if chekPathRequest(w, r, "/status") == false {
		return
	}

	checkedStatusQueues(w, r, ms.db)

}

func (ms *MServer) pingService(w http.ResponseWriter, r *http.Request) {

	if chekPathRequest(w, r, "/ping") == false {
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("pong"))

}

func (ms *MServer) CreateTables(w http.ResponseWriter, r *http.Request) {

	if chekPathRequest(w, r, "/crtab") == false {
		return
	}

	err := createTablesDB(ms.db)
	if err != nil {
		w.WriteHeader(501)
	} else {
		w.WriteHeader(200)
	}

}

func chekPathRequest(w http.ResponseWriter, r *http.Request, cPath string) bool {

	if r.URL.Path != cPath {
		http.NotFound(w, r)
		return false
	}

	return true
}
