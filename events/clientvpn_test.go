package events

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestClientVPNConnectionHandlerRequestMarshaling(t *testing.T) {
	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/clientvpn-connectionhandler-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into ClientVPNConnectionHandlerRequest
	var inputEvent ClientVPNConnectionHandlerRequest
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

func TestClientVPNConnectionHandlerRequestMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, ClientVPNConnectionHandlerRequest{})
}
