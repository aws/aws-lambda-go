package events

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events/test"
	"testing"
)

func TestAutoScalingEventMarshaling(t *testing.T) {

	// 1. read JSON from file
	inputJson := readJsonFromFile(t, "./testdata/autoscaling-event.json")

	// 2. de-serialize into Go object
	var inputEvent AutoScalingEvent
	if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. serialize to JSON
	outputJson, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 4. check result
	test.AssertJsonsEqual(t, inputJson, outputJson)
}

func TestAutoScalingMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, AutoScalingEvent{})
}
