package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

const (
	maxOpenConns    = 25
	connMaxLifetime = 15 // minutes
	maxIdleConns    = 25
	connMaxIdleTime = 10 // minutes
	maxRetries      = 10
	retryDelay      = 5 * time.Second
)

func NewPostgresDB(databaseURL string) (*DB, error) {
	var db *sqlx.DB
	var err error
	
	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", databaseURL)
		if err == nil {
			break
		}
		
		retryTime := retryDelay * time.Duration(i+1)
		log.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in %v...", 
			i+1, maxRetries, err, retryTime)
		time.Sleep(retryTime)
	}

	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres after %d attempts: %v", maxRetries, err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Minute * connMaxLifetime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(time.Minute * connMaxIdleTime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error verifying database connection: %v", err)
	}

	// Run migrations
	if err := RunMigrations(databaseURL); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DB{db}, nil
}
