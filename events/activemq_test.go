// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestActiveMQEventMarshaling(t *testing.T) {
	// 1. read JSON from file
	inputJson := test.ReadJSONFromFile(t, "./testdata/activemq-event.json")

	// 2. de-serialize into Go object
	var inputEvent ActiveMQEvent
	if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. Verify values populated into Go Object, at least one validation per data type
	assert.Equal(t, "aws:mq", inputEvent.EventSource)
	assert.Equal(t, "arn:aws:mq:us-west-2:533019413397:broker:shask-test:b-0f5b7522-2b41-4f85-a615-735a4e6d96b5", inputEvent.EventSourceARN)
	assert.Equal(t, 1, len(inputEvent.Messages))

	var message = inputEvent.Messages[0]
	assert.Equal(t, "jms/text-message", message.MessageType)
	assert.Equal(t, int64(1599863938941), message.Timestamp)
	assert.Equal(t, 1, message.DeliveryMode)
	assert.Equal(t, "testQueue", message.Destination.PhysicalName)
	assert.Equal(t, false, message.Redelivered)

	// 4. serialize to JSON
	outputJson, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 5. check result
	assert.JSONEq(t, string(inputJson), string(outputJson))
}

func TestActiveMQMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, ActiveMQEvent{})
}
