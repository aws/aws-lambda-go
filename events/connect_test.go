// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
)

func TestConnectMarshaling(t *testing.T) {

	// 1. read JSON from file
	inputJson := readJsonFromFile(t, "./testdata/connect-event.json")

	// 2. de-serialize into Go object
	var inputEvent ConnectEvent
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

func TestConnectMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, ConnectEvent{})
}
