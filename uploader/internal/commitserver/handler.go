package commitserver

import (
	"commiter/internal/storage"
	"net/http"
)

func (ms *ServerCommit) createGatewayServer(host_port string) *http.Server {

	gatewayServer := &http.Server{
		Addr:    host_port,
		Handler: ms.newMuxServe(),
	}

	return gatewayServer
}

func (ms *ServerCommit) newMuxServe() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/uploadtoquery", ms.uploadtoquery)
	mux.HandleFunc("/ping", ms.pingService)
	mux.HandleFunc("/crtab", ms.createTables)
	mux.HandleFunc("/status", ms.statusQueue)

	return mux
}
func (ms *ServerCommit) uploadtoquery(w http.ResponseWriter, r *http.Request) {

	if chekPathRequest(w, r, "/uploadtoquery") == false {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Error type Method", 405)
		return
	}
	s := storage.NewStorage(ms.DB, ms.TGBot)
	s.AddNewRequest(w, r)

}

func (ms *ServerCommit) statusQueue(w http.ResponseWriter, r *http.Request) {

	if chekPathRequest(w, r, "/status") == false {
		return
	}
	s := storage.NewStorage(ms.DB, ms.TGBot)
	s.CheckedStatusQueues(w, r)

}

func (ms *ServerCommit) pingService(w http.ResponseWriter, r *http.Request) {

	if chekPathRequest(w, r, "/ping") == false {
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("pong"))

}

func (ms *ServerCommit) createTables(w http.ResponseWriter, r *http.Request) {

	if chekPathRequest(w, r, "/crtab") == false {
		return
	}
	s := storage.NewStorage(ms.DB, ms.TGBot)
	err := s.CreateTablesDB()
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
