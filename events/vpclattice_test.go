package events

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestVPCLatticeRequestV1Marshalling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/vpclattice-v1-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent VPCLatticeRequestV1
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

func TestVPCLatticeRequestV1MalformedJson(t *testing.T) {
	test.TestMalformedJson(t, VPCLatticeRequestV1{})
}

func TestVPCLatticeRequestV2Marshalling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/vpclattice-v2-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent VPCLatticeRequestV2
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

func TestVPCLatticeRequestV2MalformedJson(t *testing.T) {
	test.TestMalformedJson(t, VPCLatticeRequestV2{})
}

func TestVPCLatticeResponse(t *testing.T) {
	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/vpclattice-response.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent VPCLatticeResponse
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
