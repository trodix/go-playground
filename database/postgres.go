package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/trodix/go-rest-api/config"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ConnectPostgres(cfg config.DBConfig) (*pgxpool.Pool, *sql.DB) {
	// Build the connection string using the config values
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %s", err)
	}

	fmt.Println("Connected to PostgreSQL")

	db, err := sql.Open("pgx", dsn)
    if err != nil {
        log.Fatalf("Unable to create sql.DB: %v\n", err)
		return nil, nil
    }

	return pool, db
}

func MigrateDatabase(db *sql.DB) {
    driver, err := pgx.WithInstance(db, &pgx.Config{})
    if err != nil {
        log.Fatalf("failed to create driver: %v", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://database/migrations", // Path to your migrations
        "pgx",              // Database driver name
        driver,
    )
    if err != nil {
        log.Fatalf("failed to create migrate instance: %v", err)
    }

    if err := m.Up(); err != nil {
        if err != migrate.ErrNoChange {
            log.Fatalf("failed to apply migrations: %v", err)
        }
    }

    log.Println("Migrations applied successfully")
}
