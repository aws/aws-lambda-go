package test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// nolint: stylecheck
func AssertJsonFile(t *testing.T, file string, o interface{}) {
	inputJSON, err := os.ReadFile(file)
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}
	AssertJsonBytes(t, inputJSON, o)
}

// nolint: stylecheck
func AssertJsonBytes(t *testing.T, inputJSON []byte, o interface{}) {
	// de-serialize
	if err := json.Unmarshal(inputJSON, o); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(o)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}
