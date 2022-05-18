package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestLexEventMarshaling(t *testing.T) {
	inputJSON := test.ReadJSONFromFile(t, "./testdata/lex-event.json")

	var inputEvent LexEvent
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestLexResponseMarshaling(t *testing.T) {
	inputJSON := test.ReadJSONFromFile(t, "./testdata/lex-response.json")

	var inputEvent LexResponse
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestLexMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, LexEvent{})
}

func TestLexResponseMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, LexResponse{})
}
