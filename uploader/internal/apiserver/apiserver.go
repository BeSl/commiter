package apiserver

import (
	"commiter/internal/storage"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

type ServerAPI struct {
	DB    *sqlx.DB
	TGBot *tgbotapi.BotAPI
}

func NewServerAPI(db *sqlx.DB, bot *tgbotapi.BotAPI) *ServerAPI {
	return &ServerAPI{
		DB:    db,
		TGBot: bot,
	}
}

func (sa *ServerAPI) CreateGatewayServer(host_port string) *http.Server {

	gatewayServer := &http.Server{
		Addr:    host_port,
		Handler: sa.newMuxServe(),
	}

	return gatewayServer
}

func (sa *ServerAPI) newMuxServe() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingService)
	mux.HandleFunc("/uploadtoquery", sa.uploadtoquery)
	mux.HandleFunc("/status", sa.statusQueue)

	return mux
}

func (sa *ServerAPI) uploadtoquery(w http.ResponseWriter, r *http.Request) {

	if chekPathRequest(w, r, "/uploadtoquery") == false {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Error type Method", 405)
		return
	}
	s := storage.NewStorage(sa.DB, sa.TGBot, nil)
	s.AddNewRequest(w, r)

}

func (sa *ServerAPI) statusQueue(w http.ResponseWriter, r *http.Request) {

	if chekPathRequest(w, r, "/status") == false {
		return
	}
	s := storage.NewStorage(sa.DB, sa.TGBot, nil)
	s.CheckedStatusQueues(w, r)

}

func pingService(w http.ResponseWriter, r *http.Request) {

	if chekPathRequest(w, r, "/ping") == false {
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("pong"))

}

func chekPathRequest(w http.ResponseWriter, r *http.Request, cPath string) bool {

	if r.URL.Path != cPath {
		http.NotFound(w, r)
		return false
	}

	return true
}
