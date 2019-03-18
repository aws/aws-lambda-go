package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestSESEventMarshaling(t *testing.T) {
	inputJSON := test.ReadJSONFromFile(t, "./testdata/ses-event.json")

	var inputEvent SimpleEmailEvent
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestSESMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, SimpleEmailEvent{})
}
