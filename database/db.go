package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

type SSLMODE string

const (
	SSLMODEDisables SSLMODE = "disable"
)

func ConnectAndMigrate(host, port, dbname, user, password string, sslmode SSLMODE) error {
	connectionStr := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		host, port, dbname, user, password, sslmode,
	)
	var err error
	DB, err = sqlx.Open("postgres", connectionStr)
	if err != nil {
		return err
	}

	//	checking connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	fmt.Println("Database connected successfully")
	return migrationUp(DB)
}

func migrationUp(db *sqlx.DB) error {
	fmt.Println("Starting database migrations...")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}
