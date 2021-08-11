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

	assert.Equal(t, inputEvent.BootstrapServers, "b-2.demo-cluster-1.a1bcde.c1.kafka.us-east-1.amazonaws.com:9092,b-1.demo-cluster-1.a1bcde.c1.kafka.us-east-1.amazonaws.com:9092")
	assert.Equal(t, inputEvent.EventSource, "aws:kafka")
	assert.Equal(t, inputEvent.EventSourceARN, "arn:aws:kafka:us-west-2:012345678901:cluster/ExampleMSKCluster/e9f754c6-d29a-4430-a7db-958a19fd2c54-4")
	for _, records := range inputEvent.Records {
		for _, record := range records {
			utc := record.Timestamp.UTC()
			assert.Equal(t, 2020, utc.Year())
			assert.Equal(t, record.Key, "OGQ1NTk2YjQtMTgxMy00MjM4LWIyNGItNmRhZDhlM2QxYzBj")
			assert.Equal(t, record.Value, "OGQ1NTk2YjQtMTgxMy00MjM4LWIyNGItNmRhZDhlM2QxYzBj")

			for _, header := range record.Headers {
				for key, value := range header {
					assert.Equal(t, key, "headerKey")
					var headerValue string = string(value)
					assert.Equal(t, headerValue, "headerValue")
				}
			}
		}
	}
}

func TestKafkaMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, KafkaEvent{})
}
