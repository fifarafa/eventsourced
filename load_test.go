package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadEmptyStream(t *testing.T) {
	streamID := uuid.New()
	events, err := loadStream(streamID)
	if err != nil {
		t.Errorf("failed to load stream: %v", err)
	}
	assert.Empty(t, events)
}

func TestLoadStreamWithSingleEvent(t *testing.T) {
	// given
	if err := setup(); err != nil {
		t.Errorf("setup failed: %v", err)
	}

	db, err := openMySQLConnection()
	if err != nil {
		t.Errorf("failed to open database: %v", err)
	}

	rawMsg := json.RawMessage(`{"key": "value"}`)

	streamID := uuid.New()
	if err := appendSingleEvent(db, streamID, "invoice", rawMsg, -1); err != nil {
		t.Errorf("failed to append event: %v", err)
	}

	// when & then
	events, err := loadStream(streamID)
	if err != nil {
		t.Errorf("failed to load stream: %v", err)
	}
	assert.Len(t, events, 1)
}

func TestLoadStreamWithMultipleEvents(t *testing.T) {

}

func TestLoadStreamWithHugeAmountOfEvents(t *testing.T) {

}
