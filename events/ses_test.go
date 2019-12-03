package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestSESEventMarshaling(t *testing.T) {

	tests := []struct {
		file string
	}{
		{"./testdata/ses-lambda-event.json"},
		{"./testdata/ses-s3-event.json"},
		{"./testdata/ses-sns-event.json"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			inputJSON := test.ReadJSONFromFile(t, tc.file)

			var inputEvent SimpleEmailEvent
			if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
				t.Errorf("could not unmarshal event. details: %v", err)
			}

			outputJSON, err := json.Marshal(inputEvent)
			if err != nil {
				t.Errorf("could not marshal event. details: %v", err)
			}

			assert.JSONEq(t, string(inputJSON), string(outputJSON))
		})
	}
}

func TestSESMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, SimpleEmailEvent{})
}
