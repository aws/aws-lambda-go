// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package events

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestConfigEventMarshaling(t *testing.T) {
	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/config-event.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent ConfigEvent
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

func TestConfigMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, ConfigEvent{})
}
