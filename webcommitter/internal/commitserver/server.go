package commitserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"webcommitter/internal/apiserver"
	"webcommitter/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ServerCommit struct {
	DB    *gorm.DB
	TGBot *tgbotapi.BotAPI
}

func New(db *gorm.DB, tgbot *tgbotapi.BotAPI) *ServerCommit {
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

	gtServer.Migrate()

	gatewayServer := gtServer.EchoServer()

	go func() {
		log.Info().Msgf("Gateway server is running on %s", gatewayAddr)
		if err := gatewayServer.Start(gatewayAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			gatewayServer.Logger.Fatal(err)
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

	return nil
}
