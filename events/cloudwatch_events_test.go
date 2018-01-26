package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
)

func TestCloudwatchScheduledEventIdempotency(t *testing.T) {
	inputJSON := []byte(
		"{\"version\":\"0\",\"id\":\"890abcde-f123-4567-890a-bcdef1234567\"," +
			"\"detail-type\":\"Scheduled Event\",\"source\":\"aws.events\"," +
			"\"account\":\"123456789012\",\"time\":\"2016-12-30T18:44:49Z\"," +
			"\"region\":\"us-east-1\"," +
			"\"resources\":[\"arn:aws:events:us-east-1:123456789012:rule/SampleRule\"]," +
			"\"detail\":{}}")

	var inputEvent CloudWatchEvent
	err := json.Unmarshal(inputJSON, &inputEvent)
	if err != nil {
		t.Errorf("Could not unmarshal scheduled event: %v", err)
	}

	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("Could not marshal scheduled event: %v", err)
	}

	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestCloudwatchScheduledEventRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CloudWatchEvent{})
}
