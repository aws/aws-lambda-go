// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestSnsEventMarshaling(t *testing.T) {

	// 1. read JSON from file
	inputJson := test.ReadJSONFromFile(t, "./testdata/sns-event.json")

	// 2. de-serialize into Go object
	var inputEvent SNSEvent
	if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. serialize to JSON
	outputJson, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 4. check result
	assert.JSONEq(t, string(inputJson), string(outputJson))
}

func TestSnsMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, SNSEvent{})
}

func TestCloudWatchAlarmSNSPayloadMarshaling(t *testing.T) {
	// 1. read JSON from file
	inputJson := test.ReadJSONFromFile(t, "./testdata/cloudwatch-alarm-sns-payload-single-metric.json")

	// 2. de-serialize into Go object
	var inputEvent SNSEvent
	if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3.retrieve message from the Go object
	var message = inputEvent.Records[0].SNS.Message

	// 4. de-serialize message into Go object
	var inputCloudWatchPayload CloudWatchAlarmSNSPayload
	if err := json.Unmarshal([]byte(message), &inputCloudWatchPayload); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 5. serialize message to JSON
	outputJson, err := json.Marshal(inputCloudWatchPayload)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 4. check result
	assert.JSONEq(t, string(message), string(outputJson))
}
