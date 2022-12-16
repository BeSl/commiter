package commitserver

import (
	"commiter/internal/apiserver"
	"commiter/internal/comittworker"

	"commiter/internal/config"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type ServerCommit struct {
	DB    *sqlx.DB
	TGBot *tgbotapi.BotAPI
}

func New(db *sqlx.DB, tgbot *tgbotapi.BotAPI) *ServerCommit {
	return &ServerCommit{
		TGBot: tgbot,
		DB:    db,
	}
}

func (ls *ServerCommit) Start(cfg *config.Config) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gatewayAddr := fmt.Sprintf("%s:%v", cfg.Rest.Host, cfg.Rest.Port)
	gtServer := apiserver.New(ls.DB, ls.TGBot)

	gatewayServer := gtServer.EchoServer()

	cm := comittworker.NewCommitCreator(ls.DB, ls.TGBot, &cfg.Gitlab)

	go func() {
		log.Info().Msgf("Gateway server is running on %s", gatewayAddr)
		if err := gatewayServer.Start(gatewayAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			gatewayServer.Logger.Fatal(err)
			cancel()
		}
	}()

	go func() {
		log.Info().Msgf("Worker commit is runnig repo %s", cfg.Gitlab.Project_url)
		if err := cm.ListenNewTasks(); err != nil {
			log.Error().
				Err(err).
				Msg("Failed running qworker job")
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

	if err := cm.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("qworker.Shutdown")
	} else {
		log.Info().Msg("worker shut down correctly")
	}
	return nil
}
