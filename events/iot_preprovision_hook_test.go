package events

import (
	"github.com/segmentio/encoding/json"
	"io/ioutil" //nolint: staticcheck
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
)

func TestIoTPreProvisionHookRequest(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/iot-preprovision-hook-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent IoTPreProvisionHookRequest
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestIoTPreProvisionHookRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, IoTPreProvisionHookRequest{})
}

func TestIoTPreProvisionHookResponseMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/iot-preprovision-hook-response.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent IoTPreProvisionHookResponse
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestIoTPreProvisionHookResponseMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, IoTPreProvisionHookResponse{})
}
