package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"kasirApi/migrations"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// getMigrateInstance creates a migrate.Migrate instance using embed.FS
func getMigrateInstance(db *sql.DB, driverName string) (*migrate.Migrate, error) {
	var driver database.Driver
	var err error

	switch driverName {
	case "mysql":
		driver, err = mysql.WithInstance(db, &mysql.Config{})
	case "postgres", "postgresql":
		driver, err = postgres.WithInstance(db, &postgres.Config{})
	default:
		driver, err = mysql.WithInstance(db, &mysql.Config{})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create %s migrate driver: %w", driverName, err)
	}

	sourceDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize iofs migration source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, driverName, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrate instance: %w", err)
	}

	return m, nil
}

// RunAutoMigrate executes all pending up migrations
func RunAutoMigrate(db *sql.DB, driverName string) error {
	m, err := getMigrateInstance(db, driverName)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration UP failed: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("[MIGRATE] Database schema is up-to-date (no changes)")
	} else {
		log.Println("[MIGRATE] Database migration UP executed successfully")
	}

	return nil
}

// RollbackMigration rolls back N migration steps
func RollbackMigration(db *sql.DB, driverName string, steps int) error {
	m, err := getMigrateInstance(db, driverName)
	if err != nil {
		return err
	}

	if err := m.Steps(-steps); err != nil {
		return fmt.Errorf("migration DOWN failed: %w", err)
	}

	log.Printf("[MIGRATE] Rolled back %d migration step(s)\n", steps)
	return nil
}
