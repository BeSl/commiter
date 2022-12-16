package main

import (
	"commiter/internal/app"
	"commiter/internal/config"

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

	err := app.New(&cfg).Start()
	if err != nil {
		log.Fatal().Err(err).Msg("App init error")
	}

}
