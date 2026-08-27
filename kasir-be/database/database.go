package database

import (
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// InitDB initializes database connection with connection pooling
func InitDB(connectionString string, driverOpt ...string) (*sql.DB, error) {
	driver := "mysql"
	if len(driverOpt) > 0 && driverOpt[0] != "" {
		driver = driverOpt[0]
	} else if strings.HasPrefix(connectionString, "postgres://") || strings.HasPrefix(connectionString, "postgresql://") {
		driver = "postgres"
	}

	db, err := sql.Open(driver, connectionString)
	if err != nil {
		return nil, err
	}

	// Connection Pool Settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	// Test Connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Printf("[DATABASE] Connected successfully via driver: %s\n", driver)
	return db, nil
}