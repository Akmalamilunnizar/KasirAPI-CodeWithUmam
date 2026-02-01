package database

import (
	"database/sql"
	"log"
	
	"github.com/lib/pq"
)

func initDB(connectionString string) (*sql.DB, error)   {
	// Open Database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Test Connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successf")
	return db, nil
}