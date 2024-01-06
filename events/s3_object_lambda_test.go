package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestS3ObjectLambdaEventMarshaling(t *testing.T) {
	tests := []struct {
		file string
	}{
		{"./testdata/s3-object-lambda-event-get-object-iam.json"},
		{"./testdata/s3-object-lambda-event-get-object-assumed-role.json"},
		{"./testdata/s3-object-lambda-event-head-object-iam.json"},
		{"./testdata/s3-object-lambda-event-list-objects-iam.json"},
		{"./testdata/s3-object-lambda-event-list-objects-v2-iam.json"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			inputJSON := test.ReadJSONFromFile(t, tc.file)

			var inputEvent S3ObjectLambdaEvent
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

func TestS3ObjectLambdaMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, S3ObjectLambdaEvent{})
}
