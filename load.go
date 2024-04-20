package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

const (
	loadStreamSQL = `
	SELECT data FROM events WHERE stream_id = ? ORDER BY created DESC`
)

func loadStream(db *sql.DB, streamID uuid.UUID) ([]json.RawMessage, error) {
	rows, err := db.Query(loadStreamSQL, streamID[:])
	if err != nil {
		return nil, fmt.Errorf("querying events from stream: %w", err)
	}
	defer rows.Close()
	var result []json.RawMessage
	for rows.Next() {
		var data json.RawMessage
		err := rows.Scan(&data)
		if err != nil {
			return nil, fmt.Errorf("scanning event data: %w", err)
		}
		result = append(result, data)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("iterating over rows: %w", err)
	}
	return result, nil
}
