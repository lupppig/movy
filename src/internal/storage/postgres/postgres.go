package postgres

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lupppig/movy/internal/logger"
)





func ConnectPostgreDB(dsn string, logger *logger.Logger) (*sql.DB, error){
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to postgres database")
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	err = db.Ping()
	if err != nil {
		logger.Error().Err(err).Msg("failed to ping postgres database")
		return nil, err
	}
	return db, err
}