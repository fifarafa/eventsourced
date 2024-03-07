package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	if err := setup(); err != nil {
		log.Fatalf("setup failed: %v", err)
	} else {
		log.Println("setup finished successfully!")
	}
}

func setup() error {
	db, err := openMySQLConnection()
	if err != nil {
		return fmt.Errorf("open mysql connection: %w", err)
	}
	defer db.Close()

	if err := setupDatabase(db); err != nil {
		return fmt.Errorf("setup database: %w", err)
	}

	return nil
}

func setupDatabase(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	if err := initTables(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("rollback transaction: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func openMySQLConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "local:local@tcp(127.0.0.1:3306)/local")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func initTables(tx *sql.Tx) error {
	if err := createStreamTable(tx); err != nil {
		return fmt.Errorf("create stream table: %w", err)
	}
	if err := createEventsTable(tx); err != nil {
		return fmt.Errorf("create events table: %w", err)
	}
	return nil
}

func createEventsTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
  		CREATE TABLE IF NOT EXISTS stream(
  		    id BINARY(16) DEFAULT (UUID_TO_BIN(UUID())) PRIMARY KEY,
  		    type VARCHAR(255) NOT NULL,
  		    version BIGINT NOT NULL
  		)`,
	)
	if err != nil {
		return fmt.Errorf("exec create events table: %w", err)
	}
	return nil
}

func createStreamTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
  		CREATE TABLE IF NOT EXISTS stream(
  		    id BINARY(16) DEFAULT (UUID_TO_BIN(UUID())) PRIMARY KEY,
  		    type VARCHAR(255) NOT NULL,
  		    version BIGINT NOT NULL
  		)`,
	)
	if err != nil {
		return fmt.Errorf("exec create stream table: %w", err)
	}
	return nil
}
