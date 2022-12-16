package app

import (
	"commiter/internal/commitserver"
	"commiter/internal/config"
	"commiter/internal/database"
	"commiter/internal/executor"
	"commiter/internal/storage"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type App struct {
	Config *config.Config
}

func New(cfg *config.Config) *App {
	return &App{
		Config: cfg,
	}
}

func (a *App) Start() error {

	exC := executor.New()
	err := exC.CloneRepo(&a.Config.Gitlab)

	if a.Config.Project.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	err = exC.Check_env()
	if err != nil {
		log.Fatal().Err(err).Msg("the environment is not initialized")
	}

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		a.Config.Database.Host,
		a.Config.Database.Port,
		a.Config.Database.User,
		a.Config.Database.Password,
		a.Config.Database.Name,
		a.Config.Database.SslMode,
	)
	db, err := database.NewPostgres(dsn, a.Config.Database.Driver)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed init postgres")
	}
	defer db.Close()

	tgbot, err := tgbotapi.NewBotAPI(a.Config.Telegramm.Token)
	if err != nil {
		log.Fatal().Err(err).Msg("Telegramm bot initial error")
	}

	if err := commitserver.New(db, tgbot).Start(a.Config); err != nil {
		log.Error().Err(err).Msg("Failed creating http server")
		return err
	}
	//Проверить таблицы в БД
	s := storage.NewStorage(db, &a.Config.Gitlab)
	s.CreateTablesDB()

	return nil
}
