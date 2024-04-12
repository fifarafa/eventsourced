package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
)

const (
	minimalSafeIsolationLevel = "READ COMMITTED"
	initialStreamVersion      = -1
)

func appendSingleEvent(db *sql.DB, streamID uuid.UUID, streamType string, event json.RawMessage, providedExpectedVersion int64) error {
	if err := setTransactionIsolationLevel(db, minimalSafeIsolationLevel); err != nil {
		return fmt.Errorf("set transaction isolation level: %w", err)
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Fatalf("rollback transaction: %v", err)
		}
	}()
	if err := appendSingleEventInner(tx, streamID, streamType, event, providedExpectedVersion); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

// TODO add support for multiple events in a stream
// TODO write benchmarks to see if tx.Prepare is faster than tx.Exec for multiple events
func appendSingleEventInner(tx *sql.Tx, streamID uuid.UUID, streamType string, event json.RawMessage, providedExpectedVersion int64) error {
	doesExist, err := streamExists(tx, streamID)
	if err != nil {
		return fmt.Errorf("checking if stream exists: %w", err)
	}
	// what if now some other transaction creates the stream?
	// this conditional insertions will work but only for commited data
	if !doesExist {
		_, err := createStream(tx, streamID, streamType)
		if err != nil {
			return fmt.Errorf("create stream: %w", err)
		}
	}

	if err := insertEvent(tx, streamID, providedExpectedVersion, event); err != nil {
		return fmt.Errorf("insert event: %w", err)
	}

	if err := incrementStreamVersion(tx, streamID, providedExpectedVersion); err != nil {
		return fmt.Errorf("conditional increment stream version: %w", err)
	}

	return nil
}

func streamExists(tx *sql.Tx, streamID uuid.UUID) (bool, error) {
	stmt, err := tx.Prepare(getStreamVersionSQL)
	if err != nil {
		return false, fmt.Errorf("prepare select stream version: %w", err)
	}
	var version int64
	err = stmt.QueryRow(streamID[:]).Scan(&version)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("get stream version: %w", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Fatalf("close statement: %v", err)
		}
	}()
	return true, nil
}

func createStream(tx *sql.Tx, streamID uuid.UUID, streamType string) (uuid.UUID, error) {
	stmt, err := tx.Prepare(insertStreamSQL)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("exec insert stream: %w", err)
	}
	res, err := stmt.Exec(
		streamID[:], streamType, initialStreamVersion,
		streamID[:],
	)
	//TODO check what error is returned here is stream already exists
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
	log.Print(streamID.String())
	return streamID, nil
}

func incrementStreamVersion(tx *sql.Tx, id uuid.UUID, version int64) error {
	stmt, err := tx.Prepare(incrementStreamVersionSQL)
	if err != nil {
		return fmt.Errorf("prepare update stream version: %w", err)
	}
	res, err := stmt.Exec(version+1, id[:], version)
	if err != nil {
		return fmt.Errorf("exec update stream version: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", rows)
	}
	return nil
}

// providedExpectedVersion should be equal to stream version saved in the db
// because it means that decision is being made on the latest state
// if it's different (smaller or bigger), the whole operation should be rejected
func insertEvent(tx *sql.Tx, streamID uuid.UUID, providedExpectedVersion int64, event json.RawMessage) error {
	stmt, err := tx.Prepare(insertEventSQL)
	if err != nil {
		return fmt.Errorf("prepare insert event: %w", err)
	}

	res, err := stmt.Exec(
		streamID[:], providedExpectedVersion+1, event, "test",
		streamID[:], providedExpectedVersion)
	if err != nil {
		return fmt.Errorf("exec insert event: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", rows)
	}
	return nil
}

func setTransactionIsolationLevel(db *sql.DB, level string) error {
	_, err := db.Exec("SET TRANSACTION ISOLATION LEVEL " + level)
	if err != nil {
		return fmt.Errorf("db exec: %w", err)
	}
	return nil
}
