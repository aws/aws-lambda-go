// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package test

import (
	"path/filepath"
	"testing"
)

type sampleEvent struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// unmarshalableEvent unmarshals cleanly from an empty JSON object but cannot be
// marshaled back to JSON because channels are an unsupported type. It is used to
// exercise the marshal-error branch of AssertJsonBytes.
type unmarshalableEvent struct {
	C chan int `json:"c"`
}

func TestAssertJsonBytes(t *testing.T) {
	// Round-trips valid JSON through Unmarshal/Marshal and asserts equality.
	inputJSON := []byte(`{"name":"test","age":5}`)
	var o sampleEvent
	AssertJsonBytes(t, inputJSON, &o)

	if o.Name != "test" || o.Age != 5 {
		t.Errorf("AssertJsonBytes did not deserialize as expected: %+v", o)
	}
}

func TestAssertJsonBytesUnmarshalError(t *testing.T) {
	// Malformed JSON should trigger the unmarshal-error branch.
	mock := &mockTB{TB: t}
	target := make(map[string]interface{})
	AssertJsonBytes(mock, []byte(`{bad json`), &target)

	if !containsError(mock.errors, "could not unmarshal event. details: %v") {
		t.Errorf("expected unmarshal error, got: %v", mock.errors)
	}
}

func TestAssertJsonBytesMarshalError(t *testing.T) {
	// A value that unmarshals from {} but cannot be marshaled should trigger the
	// marshal-error branch.
	mock := &mockTB{TB: t}
	o := &unmarshalableEvent{}
	AssertJsonBytes(mock, []byte(`{}`), o)

	if !containsError(mock.errors, "could not marshal event. details: %v") {
		t.Errorf("expected marshal error, got: %v", mock.errors)
	}
}

func TestAssertJsonFile(t *testing.T) {
	// Reads a JSON fixture and asserts it round-trips cleanly.
	var o sampleEvent
	AssertJsonFile(t, filepath.Join("testdata", "event.json"), &o)

	if o.Name != "test" || o.Age != 5 {
		t.Errorf("AssertJsonFile did not deserialize as expected: %+v", o)
	}
}

func TestAssertJsonFileMissing(t *testing.T) {
	// A missing file should trigger the open-error branch.
	mock := &mockTB{TB: t}
	var o sampleEvent
	AssertJsonFile(mock, filepath.Join("testdata", "does-not-exist.json"), &o)

	if !containsError(mock.errors, "could not open test file. details: %v") {
		t.Errorf("expected open error, got: %v", mock.errors)
	}
}

// containsError reports whether the given format string appears in errs.
func containsError(errs []string, format string) bool {
	for _, e := range errs {
		if e == format {
			return true
		}
	}
	return false
}
