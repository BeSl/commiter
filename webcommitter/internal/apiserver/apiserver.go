package apiserver

import (
	"fmt"
	"net/http"
	"webcommitter/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type ServerAPI struct {
	DB    *gorm.DB
	TGBot *tgbotapi.BotAPI
}

func New(db *gorm.DB, bot *tgbotapi.BotAPI) *ServerAPI {
	return &ServerAPI{
		DB:    db,
		TGBot: bot,
	}
}

func (sa *ServerAPI) EchoServer() *echo.Echo {

	e := echo.New()
	e.Use(Process)

	gAPI := e.Group("/api/v1")
	gAPI.GET("/ping", pingService)
	gAPI.GET("/status", sa.statusQueueCommit)

	gAPI.GET("/lastDataCommit", sa.lastDataCommit)

	gAPI.POST("/newcommit", sa.newCommit)
	gAPI.POST("/fixcommit/:id", sa.fixCommit)

	return e
}

func (sa *ServerAPI) newCommit(c echo.Context) error {

	s := storage.NewStorage(sa.DB, nil)
	msg, status, err := s.AddNewRequest(c.Response().Writer, c.Request())
	if err != nil {
		return c.String(status, fmt.Sprintf("%s . Error Description %s", msg, err.Error()))
	}
	return c.String(status, msg)
}

func (sa *ServerAPI) statusQueueCommit(c echo.Context) error {

	s := storage.NewStorage(sa.DB, nil)
	count, err := s.CheckedStatusQueues(c.Response().Writer, c.Request())
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, count)
}

func pingService(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (sa *ServerAPI) lastDataCommit(c echo.Context) error {
	return nil
}

func (sa *ServerAPI) fixCommit(c echo.Context) error {
	id := c.Param("id")
	fmt.Printf("%s", id)
	return nil
}
