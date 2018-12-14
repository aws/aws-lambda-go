package events

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestAlbLambdaTargetRequestMarshaling(t *testing.T) {
	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/alb-lambda-target-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent LambdaTargetGroupRequest
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

func TestAlbLambdaTargetRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, LambdaTargetGroupRequest{})
}

func TestAlbLambdaTargetResponseMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/alb-lambda-target-response.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent LambdaTargetGroupResponse
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
