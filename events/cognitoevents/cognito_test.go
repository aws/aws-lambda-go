// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package cognitoevents

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
)

func TestCognitoEventMarshaling(t *testing.T) {

	// read json from file
	inputJson, err := ioutil.ReadFile("./test-files/cognito-event.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into CognitoEvent
	var inputEvent CognitoEvent
	if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJson, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	test.AssertJsonsEqual(t, inputJson, outputJson)
}

func TestCognitoMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CognitoEvent{})
}
