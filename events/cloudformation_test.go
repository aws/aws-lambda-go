package events

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
)

func TestCloudformationCustomResourceRequestIdempotency(t *testing.T) {
	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/cloudformation-custom-resource-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	var inputEvent CloudformationCustomResourceRequest
	err = json.Unmarshal(inputJSON, &inputEvent)
	if err != nil {
		t.Errorf("Could not unmarshal scheduled event: %v", err)
	}

	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("Could not marshal scheduled event: %v", err)
	}

	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestCloudformationCustomResourceRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CloudformationCustomResourceRequest{})
}

func TestCloudformationCustomResourceResponseIdempotency(t *testing.T) {
	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/cloudformation-custom-resource-response.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	var inputEvent CloudformationCustomResourceResponse
	err = json.Unmarshal(inputJSON, &inputEvent)
	if err != nil {
		t.Errorf("Could not unmarshal scheduled event: %v", err)
	}

	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("Could not marshal scheduled event: %v", err)
	}

	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestCloudformationCustomResourceResponseMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CloudformationCustomResourceResponse{})
}
