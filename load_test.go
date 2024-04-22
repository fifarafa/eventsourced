package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestLoadEmptyStream(t *testing.T) {
	db, err := openMySQLConnection()
	if err != nil {
		t.Errorf("failed to open database: %v", err)
	}
	streamID := uuid.New()
	events, err := loadStream(db, streamID)
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
	events, err := loadStream(db, streamID)
	if err != nil {
		t.Errorf("failed to load stream: %v", err)
	}
	assert.Len(t, events, 1)
}

func TestLoadStreamWithMultipleEvents(t *testing.T) {
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

	if err := appendSingleEvent(db, streamID, "invoice", rawMsg, 0); err != nil {
		t.Errorf("failed to append event: %v", err)
	}

	if err := appendSingleEvent(db, streamID, "invoice", rawMsg, 1); err != nil {
		t.Errorf("failed to append event: %v", err)
	}

	// when & then
	events, err := loadStream(db, streamID)
	if err != nil {
		t.Errorf("failed to load stream: %v", err)
	}
	assert.Len(t, events, 3)
}

func TestLoadStreamWithBigAmountOfEvents(t *testing.T) {
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

	for i := 0; i < 1000; i++ {
		if err := appendSingleEvent(db, streamID, "invoice", rawMsg, int64(i-1)); err != nil {
			t.Errorf("failed to append event: %v", err)
		}
	}
	// when & then
	events, err := loadStream(db, streamID)
	if err != nil {
		t.Errorf("failed to load stream: %v", err)
	}
	assert.Len(t, events, 1000)
}

func TestLoadStreamInConcurrentEnvironment(t *testing.T) {
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

	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			if err := appendSingleEvent(db, streamID, "invoice", rawMsg, int64(i-1)); err != nil {
				t.Errorf("failed to append event: %v", err)
			}
		}()
	}
	wg.Wait()
	// when & then
	events, err := loadStream(db, streamID)
	if err != nil {
		t.Errorf("failed to load stream: %v", err)
	}
	assert.Len(t, events, 1000)
}
