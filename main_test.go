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
	streamID := uuid.New()
	if err := appendSingleEvent(db, streamID, "invoice", rawMsg, -1); err != nil {
		t.Errorf("failed to append event: %v", err)
	}

	events, err := loadStream(streamID)
	if err != nil {
		t.Errorf("failed to load stream: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

func TestWhenStreamWithProvidedVersionAlreadyExists(t *testing.T) {
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
	streamID := uuid.New()
	if err := appendSingleEvent(db, streamID, "invoice", rawMsg, -1); err != nil {
		t.Errorf("failed to append event: %v", err)
	}

	err = appendSingleEvent(db, streamID, "invoice", rawMsg, -1)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
