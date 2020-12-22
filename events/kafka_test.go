// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestKafkaEventMarshaling(t *testing.T) {

	// 1. read JSON from file
	inputJson := test.ReadJSONFromFile(t, "./testdata/kafka-event.json")

	// 2. de-serialize into Go object
	var inputEvent KafkaEvent
	if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	for _, records := range inputEvent.Records {
		for _, record := range records {
			utc := record.Timestamp.UTC()
			assert.Equal(t, 2020, utc.Year())
		}
	}

	// 3. serialize to JSON
	outputJson, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 4. check result
	assert.JSONEq(t, string(inputJson), string(outputJson))
}

func TestKafkaMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, KafkaEvent{})
}
