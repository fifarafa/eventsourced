package main

import (
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

	// when & then
	if err := appendSingleEvent(db, uuid.New(), nil, 0); err != nil {
		t.Errorf("failed to append event: %v", err)
	}
}
