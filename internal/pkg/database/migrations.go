package database

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB) error {
	log.Println("Running database migrations...")
	
	dir, err := os.Getwd()
	if err != nil {
		 log.Printf("Error getting current directory: %v", err)
	}
	log.Printf("Current working directory: %s", dir)

	// Add this to check migration files
	migrationPath := "pkg/database/migrations"
	files, err := os.ReadDir(migrationPath)
	if err != nil {
		 log.Printf("Error reading migrations directory: %v", err)
	} else {
		 log.Println("Migration files found:")
		 for _, file := range files {
			  log.Printf("- %s", file.Name())
		 }
	}

	sqlDB := db.DB
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		 log.Printf("Driver error: %v", err)
		 return fmt.Errorf("could not create the postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		 "file://internal/pkg/database/migrations",
		 "postgres", driver)
	if err != nil {
		 log.Printf("Migration initialization error: %v", err)
		 return fmt.Errorf("migration initialization failed: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		 log.Printf("Error getting migration version: %v", err)
	} else {
		 log.Printf("Current migration version: %d, Dirty: %v", version, dirty)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		 log.Printf("Migration error: %v", err)
		 return fmt.Errorf("migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}