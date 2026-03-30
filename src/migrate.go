package main

import (
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(dsn string) error {
	// DEBUG: List all embedded files


	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
    if err == migrate.ErrNoChange {
        log.Println("✅ Database is already up to date (Version 1 detected)")
    } else {
        log.Printf("❌ Migration failed: %v", err)
        return err
    }
	} else {
		log.Println("🚀 SUCCESS: New migrations were applied to the database!")
	}
	return nil
}
