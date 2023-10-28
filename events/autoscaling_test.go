package events

import (
	"github.com/segmentio/encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestAutoScalingEventMarshaling(t *testing.T) {

	var sampleFileList = []string{"autoscaling-event-launch-successful.json", "autoscaling-event-launch-unsuccessful.json", "autoscaling-event-lifecycle-action.json",
		"autoscaling-event-terminate-action.json", "autoscaling-event-terminate-successful.json", "autoscaling-event-terminate-unsuccessful.json"}

	// Loop over list and test each file individually //
	for _, sampleFile := range sampleFileList {

		t.Logf("Running test for %s\n", sampleFile)
		// 1. read JSON from file
		inputJSON := test.ReadJSONFromFile(t, "./testdata/"+sampleFile)

		// 2. de-serialize into Go object
		var inputEvent AutoScalingEvent
		if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
			t.Errorf("could not unmarshal event. details: %v", err)
		}
		// 3. serialize to JSON
		outputJSON, err := json.Marshal(inputEvent)
		if err != nil {
			t.Errorf("could not marshal event. details: %v", err)
		}
		// 4. check result
		assert.JSONEq(t, string(inputJSON), string(outputJSON))
	}

}

func TestAutoScalingMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, AutoScalingEvent{})
}
