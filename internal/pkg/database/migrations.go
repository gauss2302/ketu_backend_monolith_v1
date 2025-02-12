package database

import (
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations runs database migrations
func RunMigrations(dsn string) error {
	log.Printf("Starting database migrations with DSN: %s", dsn)
	
	// Print available files in embedded FS for debugging
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		log.Printf("Error reading migrations dir: %v", err)
	} else {
		log.Printf("Available migration files:")
		for _, entry := range entries {
			log.Printf("- %s", entry.Name())
		}
	}
	
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Force to the latest version
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("Warning: could not get migration version: %v", err)
	} else {
		log.Printf("Migration completed. Version: %d, Dirty: %v", version, dirty)
	}

	return nil
}
