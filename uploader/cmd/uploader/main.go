package main

import (
	"commiter/internal/commitserver"
	"commiter/internal/config"
	"commiter/internal/database"
	"commiter/internal/executor"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}
	cfg := config.GetConfigInstance()

	log.Info().
		Str("version", cfg.Project.Version).
		Str("commitHash", cfg.Project.CommitHash).
		Bool("debug", cfg.Project.Debug).
		Str("environment", cfg.Project.Environment).
		Msgf("Starting service: %s", cfg.Project.Name)

	if cfg.Project.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	exC := executor.NewExecutor()
	err := exC.Check_env()
	if err != nil {
		log.Fatal().Err(err).Msg("the environment is not initialized")
	}
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)
	db, err := database.NewPostgres(dsn, cfg.Database.Driver)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed init postgres")
	}
	defer db.Close()

	tgbot, err := tgbotapi.NewBotAPI(cfg.Telegramm.Token)
	if err != nil {
		log.Fatal().Err(err).Msg("Telegramm bot initial error")
	}

	if err := commitserver.NewlServer(db, tgbot).Start(&cfg); err != nil {
		log.Error().Err(err).Msg("Failed creating http server")
		return
	}
}
