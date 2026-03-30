package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lupppig/movy/internal/config"
	"github.com/lupppig/movy/internal/logger"
	"github.com/lupppig/movy/internal/storage/postgres"
)

//go:generate oapi-codegen --config=./internal/openapi/openapi-config/config.yaml ./internal/openapi/openapi.yaml

func main()  {
	logger :=  logger.NewLogger()
	config, err := config.LoadConfig()

	if err != nil {
		logger.Error().Err(err).Msg("failed to load config variables...")
	}
	db, err := postgres.ConnectPostgreDB(config.DATABASE_URL, logger)
	if err != nil {
		logger.Error().Err(err).Msg("connecting to database failed check DSN string or database is up and running")
		return
	}


	// run go migration
	if err := RunMigrations(config.DATABASE_URL); err != nil {
		logger.Error().Err(err).Msg("failed to perform migration")
		return
	}

	router := Router(*config, logger, db)
	port := config.APP_PORT

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        router,               
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,          // 1MB
	}

	logger.Info().Str("port", port).Msg("server is successfully booted up")
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal().Err(err).Msg(fmt.Sprintf("failed to start server in port %v", port))
	}
}