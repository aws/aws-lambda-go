package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestLexV2EventMarshaling(t *testing.T) {
	inputJSON := test.ReadJSONFromFile(t, "./testdata/lexv2-event.json")

	var inputEvent LexV2Event
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestLexV2ResponseMarshaling(t *testing.T) {
	inputJSON := test.ReadJSONFromFile(t, "./testdata/lexv2-response.json")

	var inputEvent LexV2Response
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestLexV2MarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, LexV2Event{})
}

func TestLexV2ResponseMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, LexV2Response{})
}
