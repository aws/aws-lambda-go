package events

import (
	"github.com/segmentio/encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestRabbitMQEventMarshaling(t *testing.T) {
	// 1. read JSON from file
	inputJSON := test.ReadJSONFromFile(t, "./testdata/rabbitmq-event.json")

	// 2. de-serialize into Go object
	var inputEvent RabbitMQEvent
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. Verify values populated into Go Object, at least one validation per data type
	assert.Equal(t, "aws:rmq", inputEvent.EventSource)
	assert.Equal(t, "arn:aws:mq:us-west-2:112556298976:broker:test:b-9bcfa592-423a-4942-879d-eb284b418fc8", inputEvent.EventSourceARN)
	assert.Equal(t, 1, len(inputEvent.MessagesByQueue))
	for _, messages := range inputEvent.MessagesByQueue {
		for _, message := range messages {
			assert.Equal(t, false, message.Redelivered)
			assert.Equal(t, "text/plain", message.BasicProperties.ContentType)
		}
	}

	// 4. serialize to JSON
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 5. check result
	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestRabbitMQMarshalingMalformedJSON(t *testing.T) {
	test.TestMalformedJson(t, RabbitMQEvent{})
}
