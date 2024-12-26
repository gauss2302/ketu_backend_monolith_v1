package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB) error {
	log.Println("Running database migrations...")

	// Get the underlying sql.DB object
	sqlDB := db.DB

	// Create a new postgres driver for migrations
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create the postgres driver: %v", err)
	}

	// Create a new migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/database/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("migration initialization failed: %v", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
