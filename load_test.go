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
	_ = json.RawMessage(`{"key": "value"}`)
}

func TestLoadStreamWithMultipleEvents(t *testing.T) {

}
