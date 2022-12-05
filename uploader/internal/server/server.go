package server

import (
	"commiter/internal/config"
	"commiter/internal/qworker"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type MServer struct {
	db *sqlx.DB
	//db    string
	tgBot *tgbotapi.BotAPI
}

func NewlServer(db *sqlx.DB, tgbot *tgbotapi.BotAPI) *MServer {
	//func NewlServer(db string, tgbot *tgbotapi.BotAPI) *MServer {
	return &MServer{
		tgBot: tgbot,
		db:    db,
	}
}

func (ls *MServer) Start(cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gatewayAddr := fmt.Sprintf("%s:%v", cfg.Rest.Host, cfg.Rest.Port)
	gatewayServer := ls.createGatewayServer(gatewayAddr)
	qw := qworker.NewQWorker(&cfg.Gitlab)

	go func() {
		log.Info().Msgf("Gateway server is running on %s", gatewayAddr)
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Failed running gateway server")
			cancel()
		}
	}()

	go func() {
		log.Info().Msgf("Worker commit is runnig repo %s", cfg.Gitlab.Project_url)

		if err := qw.ListenNewJob(); err != nil {
			log.Error().Err(err).Msg("Failed running qworker job")
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		log.Info().Msgf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		log.Info().Msgf("ctx.Done: %v", done)
	}

	if err := gatewayServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("gatewayServer.Shutdown")
	} else {
		log.Info().Msg("gatewayServer shut down correctly")
	}

	if err := qw.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("qworker.Shutdown")
	} else {
		log.Info().Msg("worker shut down correctly")
	}
	return nil
}
