// Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestKinesisAnalyticsOutputDeliveryEventMarshaling(t *testing.T) {
	testKinesisAnalyticsOutputMarshaling(t, KinesisAnalyticsOutputDeliveryEvent{}, "./testdata/kinesis-analytics-output-delivery-event.json")
}

func TestKinesisAnalyticsOutputDeliveryResponseMarshaling(t *testing.T) {
	testKinesisAnalyticsOutputMarshaling(t, KinesisAnalyticsOutputDeliveryResponse{}, "./testdata/kinesis-analytics-output-delivery-response.json")
}

func TestKinesisOutputDeliveryEventMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, KinesisAnalyticsOutputDeliveryEvent{})
}

func testKinesisAnalyticsOutputMarshaling(t *testing.T, inputEvent interface{}, jsonFile string) {

	// 1. read JSON from file
	inputJSON := test.ReadJSONFromFile(t, jsonFile)

	// 2. de-serialize into Go object
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. serialize to JSON
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 4. check result
	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}
