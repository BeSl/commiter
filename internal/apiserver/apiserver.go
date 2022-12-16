package apiserver

import (
	"commiter/internal/storage"
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type ServerAPI struct {
	DB    *sqlx.DB
	TGBot *tgbotapi.BotAPI
}

func New(db *sqlx.DB, bot *tgbotapi.BotAPI) *ServerAPI {
	return &ServerAPI{
		DB:    db,
		TGBot: bot,
	}
}

func (sa *ServerAPI) EchoServer(host_port string) *echo.Echo {

	e := echo.New()
	e.GET("/ping", pingService)
	e.GET("/status", sa.statusQueue)

	e.POST("/uploadtoquery", sa.uploadtoquery)

	return e
}

func (sa *ServerAPI) uploadtoquery(c echo.Context) error {

	s := storage.NewStorage(sa.DB, nil)
	msg, status, err := s.AddNewRequest(c.Response().Writer, c.Request())
	if err != nil {
		return c.String(status, fmt.Sprintf("%s . Error Description %s", msg, err.Error()))
	}
	return c.String(status, msg)
}

func (sa *ServerAPI) statusQueue(c echo.Context) error {

	s := storage.NewStorage(sa.DB, nil)
	err := s.CheckedStatusQueues(c.Response().Writer, c.Request())
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "_OK")
}

func pingService(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
