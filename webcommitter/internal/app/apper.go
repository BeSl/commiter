package app

import (
	"fmt"
	"webcommitter/internal/commitserver"
	"webcommitter/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	if a.Config.Project.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		a.Config.Database.Host,
		a.Config.Database.Port,
		a.Config.Database.User,
		a.Config.Database.Password,
		a.Config.Database.Name,
		a.Config.Database.SslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed init DB connection")
	}
	
	bot, err := tgbotapi.NewBotAPI(a.Config.Telegramm.Token)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed init telegramm bot")
	}

	err = commitserver.New(db, bot).Start(a.Config)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed init server commit")
	}

	return nil
}
