package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"ketu_backend_monolith_v1/configs"
	"time"
)

func NewPostgresDB(cfg *configs.PostgresConfig) (*sqlx.DB, error) {
	// Connect to the database
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Minute * connMaxLifetime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(time.Second * connMaxIdleTime)

	return db, nil
}

const (
	maxOpenConns    = 25 // Total number of concurrent connections
	connMaxLifetime = 15 // Minutes before connection is closed
	maxIdleConns    = 25 // Number of idle connections maintained
	connMaxIdleTime = 10 // Minutes before idle connection is closed
)
