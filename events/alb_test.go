package events

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestALBTargetRequestMarshaling(t *testing.T) {
	inputFiles := []string{
		"./testdata/alb-lambda-target-request-headers-only.json",
		"./testdata/alb-lambda-target-request-multivalue-headers.json"}

	for _, filename := range inputFiles {
		// read json from file
		inputJSON, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Errorf("could not open test file. details: %v", err)
		}

		// de-serialize into Go object
		var inputEvent ALBTargetGroupRequest
		if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
			t.Errorf("could not unmarshal event. details: %v", err)
		}

		// serialize to json
		outputJSON, err := json.Marshal(inputEvent)
		if err != nil {
			t.Errorf("could not marshal event. details: %v", err)
		}
		assert.JSONEq(t, string(inputJSON), string(outputJSON))
	}
}

func TestALBTargetRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, ALBTargetGroupRequest{})
}

func TestALBTargetResponseMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/alb-lambda-target-response.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent ALBTargetGroupResponse
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}
