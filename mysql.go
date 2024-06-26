package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

const (
	createStreamsTableSQL = `
    	CREATE TABLE IF NOT EXISTS streams(
    	    id BINARY(16) PRIMARY KEY,
    	    type VARCHAR(255) NOT NULL,
    	    version BIGINT NOT NULL
    	)`
	createEventsTableSQL = `
    	CREATE TABLE IF NOT EXISTS events(
    	    id BINARY(16) DEFAULT (UUID_TO_BIN(UUID())) PRIMARY KEY,
			stream_id BINARY(16) NOT NULL,
			data JSON NOT NULL,
    	    type VARCHAR(255) NOT NULL,
    	    version BIGINT NOT NULL,
			created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT events_stream_stream_id_fk FOREIGN KEY (stream_id) REFERENCES streams(id)
    	)`
)

func setup() error {
	db, err := openMySQLConnection()
	if err != nil {
		return fmt.Errorf("open mysql connection: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	if err := setupDatabase(db); err != nil {
		return fmt.Errorf("setup database: %w", err)
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

func setupDatabase(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Fatalf("rollback transaction: %v", err)
		}
	}()
	if err := initTables(tx); err != nil {
		return fmt.Errorf("init tables: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func initTables(tx *sql.Tx) error {
	if err := createTable(createStreamsTableSQL, tx); err != nil {
		return fmt.Errorf("create stream table: %w", err)
	}
	if err := createTable(createEventsTableSQL, tx); err != nil {
		return fmt.Errorf("create events table: %w", err)
	}
	return nil
}

func createTable(sqlStmt string, tx *sql.Tx) error {
	_, err := tx.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("exec create table: %w", err)
	}
	return nil
}
