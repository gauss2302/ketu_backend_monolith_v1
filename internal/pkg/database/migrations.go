package database

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"log"
)

func RunMigrations(db *sqlx.DB) error {
	log.Println("Running database migrations...")

	sqlDB := db.DB
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Printf("Driver error: %v", err)
		return fmt.Errorf("could not create the postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/database/migrations",
		"myapp", driver)
	if err != nil {
		log.Printf("Migration initialization error: %v", err)
		return fmt.Errorf("migration initialization failed: %v", err)
	}

	// Check and log the current migration version
	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		log.Printf("Error getting migration version: %v", err)
	} else {
		log.Printf("Current migration version: %d, Dirty: %v", version, dirty)
	}

	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Migration error: %v", err)
		return fmt.Errorf("migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
