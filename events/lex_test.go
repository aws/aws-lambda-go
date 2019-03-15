package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestLexEventMarshaling(t *testing.T) {
	tests := []struct {
		filePath string
	}{{"./testdata/lex-response.json"}, {"./testdata/lex-event.json"}}

	for _, te := range tests {
		inputJSON := test.ReadJSONFromFile(t, te.filePath)

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
}

func TestLexMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, LexEvent{})
}
