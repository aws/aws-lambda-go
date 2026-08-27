// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package test

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestReadJSONFromFile(t *testing.T) {
	// Read a known fixture and confirm the bytes deserialize as expected.
	got := ReadJSONFromFile(t, filepath.Join("testdata", "event.json"))

	var o sampleEvent
	if err := json.Unmarshal(got, &o); err != nil {
		t.Fatalf("ReadJSONFromFile returned invalid JSON: %v", err)
	}
	if o.Name != "test" || o.Age != 5 {
		t.Errorf("ReadJSONFromFile returned unexpected content: %s", got)
	}
}

func TestReadJSONFromFileMissing(t *testing.T) {
	// A missing file should trigger Errorf and return nil.
	mock := &mockTB{}
	got := ReadJSONFromFile(mock, filepath.Join("testdata", "does-not-exist.json"))

	if got != nil {
		t.Errorf("expected nil content for missing file, got: %s", got)
	}
	if len(mock.errors) != 1 {
		t.Fatalf("expected 1 error for missing file, got %d: %v", len(mock.errors), mock.errors)
	}
	if mock.errors[0] != "could not open test file. details: %v" {
		t.Errorf("unexpected error message: %s", mock.errors[0])
	}
}
