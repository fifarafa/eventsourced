package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"log"
)

const (
	createStreamsTableSQL = `
    CREATE TABLE IF NOT EXISTS streams(
        id BINARY(16) DEFAULT (UUID_TO_BIN(UUID())) PRIMARY KEY,
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
		CONSTRAINT events_stream_stream_id_fk FOREIGN KEY (stream_id) REFERENCES stream(id)
    )`
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
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
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

// TODO add support for multiple events in a stream
// TODO write benchmarks to see if tx.Prepare is faster than tx.Exec for multiple events
func appendSingleEvent(db *sql.DB, streamID uuid.UUID, event json.RawMessage, expectedVersion int64) error {
	// get stream version
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	strVer, err := getStreamVersion(tx, streamID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		newStreamID, err := createStream(tx, "test")
		if err != nil {
			return fmt.Errorf("create stream: %w", err)
		}
		streamID = newStreamID
	case err != nil:
		return fmt.Errorf("get stream version: %w", err)
	}
	log.Println("streamID", streamID, "version", strVer)
	return nil
	// check optimistic concurrency
	// append event with version = stream_version + 1
	// update stream with version = stream_version + 1
}

func getStreamVersion(tx *sql.Tx, streamID uuid.UUID) (int64, error) {
	stmt, err := tx.Prepare("SELECT version FROM stream WHERE id = (?)")
	if err != nil {
		return 0, fmt.Errorf("prepare select stream version: %w", err)
	}
	var version int64
	err = stmt.QueryRow(streamID).Scan(version)
	if err != nil {
		return 0, fmt.Errorf("query row: %w", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Fatalf("close statement: %v", err)
		}
	}()
	return version, nil
}

func createStream(tx *sql.Tx, streamType string) (uuid.UUID, error) {
	stmt, err := tx.Prepare("INSERT INTO streams (id, type, version) VALUES (?, ?, ?)")
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("exec insert stream: %w", err)
	}
	streamID := uuid.New()
	res, err := stmt.Exec(streamID[:], streamType, 0)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("exec insert stream: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("rows affected: %w", err)
	}
	if rows != 1 {
		return uuid.UUID{}, fmt.Errorf("expected 1 row affected, got %d", rows)
	}
	return streamID, nil
}
