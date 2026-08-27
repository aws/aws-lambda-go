// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package test

import (
	"encoding/json"
	"errors"
	"testing"
)

// mockTB implements testing.TB to capture Errorf calls.
type mockTB struct {
	testing.TB
	errors []string
}

func (m *mockTB) Errorf(format string, args ...interface{}) {
	m.errors = append(m.errors, format)
}

func (m *mockTB) Helper() {}

func TestMalformedJsonSyntaxError(t *testing.T) {
	// Default behavior: malformed JSON produces json.SyntaxError, no Errorf called.
	mock := &mockTB{}
	target := make(map[string]interface{})
	TestMalformedJson(mock, &target)

	if len(mock.errors) != 0 {
		t.Errorf("expected no errors, got: %v", mock.errors)
	}
}

func TestMalformedJsonUnmarshalSucceeds(t *testing.T) {
	// When unmarshal unexpectedly returns nil, both Errorf calls should fire.
	original := unmarshalJSON
	unmarshalJSON = func([]byte, interface{}) error {
		return nil
	}
	defer func() { unmarshalJSON = original }()

	mock := &mockTB{}
	TestMalformedJson(mock, nil)

	if len(mock.errors) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(mock.errors), mock.errors)
	}
	if mock.errors[0] != "unmarshal should have failed but succeeded instead" {
		t.Errorf("unexpected first error message: %s", mock.errors[0])
	}
	if mock.errors[1] != "unmarshal should have returned a json.SyntaxError" {
		t.Errorf("unexpected second error message: %s", mock.errors[1])
	}
}

func TestMalformedJsonNonSyntaxError(t *testing.T) {
	// When unmarshal returns a non-SyntaxError, the second Errorf should fire.
	original := unmarshalJSON
	unmarshalJSON = func([]byte, interface{}) error {
		return errors.New("some other error")
	}
	defer func() { unmarshalJSON = original }()

	mock := &mockTB{}
	TestMalformedJson(mock, nil)

	if len(mock.errors) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(mock.errors), mock.errors)
	}
	if mock.errors[0] != "unmarshal should have returned a json.SyntaxError" {
		t.Errorf("unexpected error message: %s", mock.errors[0])
	}
}

func TestMalformedJsonDefaultUsesStdlib(t *testing.T) {
	// Verify the default unmarshalJSON is json.Unmarshal.
	mock := &mockTB{}
	var target json.RawMessage
	TestMalformedJson(mock, &target)

	if len(mock.errors) != 0 {
		t.Errorf("expected no errors with default json.Unmarshal, got: %v", mock.errors)
	}
}
