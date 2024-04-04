package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"testing"
)

// TODO run this test using 3rd party library for setting up and tearing down the database
func TestWhenNoStreamsYet(t *testing.T) {
	// given
	if err := setup(); err != nil {
		t.Errorf("setup failed: %v", err)
	}

	db, err := openMySQLConnection()
	if err != nil {
		t.Errorf("failed to open database: %v", err)
	}

	rawMsg := json.RawMessage(`{"key": "value"}`)
	// when & then
	if err := appendSingleEvent(db, uuid.New(), rawMsg, -1); err != nil {
		t.Errorf("failed to append event: %v", err)
	}
}
