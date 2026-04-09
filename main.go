package main

import (
	"os"
	"print-agent/cmd"
	"print-agent/config"
	"print-agent/internal/adapter/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	err := config.LoadConfig("../")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to load config file")
	}

	db, err := database.NewSqliteConn()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to connect to database")
	}

	app := &cmd.App{}

	router := fiber.New()
	app.Bootstrap(router, db)

	cmd.Execute()
}
